package repository

import (
	"database/sql"

	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type IHiveosRepository interface {
	GetBlock24NotUnckle() int
	GetBlock24Uncle() int
	GetBillInfo(walletId string) RepoBillInfo
	GetBill(walletId string) *sql.Rows
	GetBalance(walletId string) float64
}

type HiveosRepository struct {
	hiveClient *gorm.DB
}

type RepoBillInfo struct {
	Balance   string
	FirstPaid float64
	FirstTime string
	TotalPaid float64
}

type Payment struct {
	firstPaid float64
	paymentId float64
}

func NewHiveosRepository() IHiveosRepository {
	db, err := mysql.Connect(config.Sequelize2DB)

	if err != nil {
		log.Panic("failed to init mysql db :", err.Error())
	}

	return &HiveosRepository{hiveClient: db}
}

func (repo *HiveosRepository) queryIntSingle(query string) int {
	var result int
	row := repo.hiveClient.Raw(query).Row()
	row.Scan(&result)
	return result
}

func (repo *HiveosRepository) GetBlock24NotUnckle() int {
	sql :=
		`SELECT count(id) as count 
		FROM blocks as b 
		WHERE is_uncle=0 and 
		b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval ` + config.PgOneDay + `)) * 1000`
	return repo.queryIntSingle(sql)
}

func (repo *HiveosRepository) GetBlock24Uncle() int {
	sql :=
		`SELECT count(id) as count 
		FROM blocks as b 
		WHERE is_uncle=1 and 
		b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval ` + config.PgOneDay + `)) * 1000`
	return repo.queryIntSingle(sql)
}

func (repo *HiveosRepository) GetBill(walletId string) *sql.Rows {
	sql :=
		`SELECT p.id, pd.paid, p.status, p.create_ts, p.tx_hash
		FROM payment_details pd
		INNER JOIN payments p ON p.id = pd.id 
		WHERE pd.miner_wallet = ? ORDER BY pd.id DESC LIMIT 30`

	rows, err := repo.hiveClient.Raw(sql, walletId).Rows()
	if err != nil {
		log.Error(err)
	}

	return rows
}

func (repo *HiveosRepository) GetBillInfo(walletId string) RepoBillInfo {
	var totalPaid float64

	err := repo.hiveClient.Raw(
		`SELECT sum(paid) AS totalPaid 
		FROM payment_details 
		WHERE miner_wallet = ?`, walletId).Scan(&totalPaid)

	if err != nil {
		log.Error(err)
	}

	var payment Payment
	var firstTime, balance string
	err = repo.hiveClient.Raw(
		`SELECT paid,payment_id FROM payment_details 
		WHERE miner_wallet = ? ORDER BY id LIMIT 1`, walletId).Scan(&payment)

	if err != nil {
		log.Error(err)
	}

	repo.hiveClient.Raw(
		`SELECT create_ts 
		FROM payments 
		WHERE id=?`, payment.paymentId).Scan(&firstTime)

	if err != nil {
		log.Error(err)
	}

	repo.hiveClient.Raw(
		`SELECT balance 
		FROM deposits 
		WHERE miner_wallet =?`, walletId).Scan(&balance)
	if err != nil {
		log.Error(err)
	}

	return RepoBillInfo{Balance: balance, FirstPaid: payment.firstPaid, FirstTime: firstTime, TotalPaid: totalPaid}

}

func (repo *HiveosRepository) GetBalance(walletId string) float64 {
	var res float64
	err := repo.hiveClient.Raw(
		`SELECT balance 
		FROM deposits 
		WHERE miner_wallet = ?`, walletId).Scan(&res)
	if err != nil {
		log.Error(err)
	}
	return res
}
