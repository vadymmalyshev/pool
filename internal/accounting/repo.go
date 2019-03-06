package accounting

import (
	"database/sql"
	"git.tor.ph/hiveon/pool/config"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type AccointingRepositorer interface {
	GetBlock24NotUnckle() int
	GetBlock24Uncle() int
	GetBillInfo(walletId string) RepoBillInfo
	GetBill(walletId string) *sql.Rows
	GetBalance(walletId string) float64
}

type AccountingRepository struct {
	db *gorm.DB
}

func NewAccountingRepository(db *gorm.DB) AccointingRepositorer {
	return &AccountingRepository{db: db}
}

func (repo *AccountingRepository) queryIntSingle(query string) int {
	var result int
	row := repo.db.Raw(query).Row()
	row.Scan(&result)
	return result
}

func (repo *AccountingRepository) GetBlock24NotUnckle() int {
	sql := "select count(id) as count from blocks as b where is_uncle=0 and b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval " +
		config.PgOneDay + ")) * 1000"
	return repo.queryIntSingle(sql)
}

func (repo *AccountingRepository) GetBlock24Uncle() int {
	sql := "select count(id) as count from blocks as b where is_uncle=1 and b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval " +
		config.PgOneDay + ")) * 1000"
	return repo.queryIntSingle(sql)
}

func (repo *AccountingRepository) GetBill(walletId string) *sql.Rows {
	rows, err := repo.db.Raw("select p.id, pd.paid, p.status, p.create_ts, p.tx_hash  from payment_details pd " +
		"inner join payments p on p.id = pd.id where pd.miner_wallet = ? order by pd.id desc limit 30", walletId).Rows()
	if err != nil {
		log.Error(err)
	}

	return rows
}

func (repo *AccountingRepository) GetBillInfo(walletId string) RepoBillInfo {
	var totalPaid float64
	err := repo.db.Raw("select sum(paid) as totalPaid from payment_details where miner_wallet = ?", walletId).Row().Scan(&totalPaid)
	if err != nil {
		log.Error(err)
	}

	var payment Payment
	var firstTime, balance string
	err1 := repo.db.Raw("select paid,payment_id from payment_details where miner_wallet = ? order by id limit 1",
		walletId).Row().Scan(&payment.firstPaid, &payment.paymentId)
	if err1 != nil {
		log.Error(err)
	}

	repo.db.Raw("select create_ts from payments where id = ?", payment.paymentId).Row().Scan(&firstTime)
	if err != nil {
		log.Error(err)
	}

	repo.db.Raw("select balance from deposits where miner_wallet = ?", walletId).Row().Scan(&balance)
	if err != nil {
		log.Error(err)
	}

	return RepoBillInfo{Balance: balance, FirstPaid: payment.firstPaid, FirstTime: firstTime, TotalPaid: totalPaid}

}

func (repo *AccountingRepository) GetBalance(walletId string) float64 {
	var res float64
	err := repo.db.Raw("select balance from deposits where miner_wallet = ?", walletId).Row().Scan(&res)
	if err != nil {
		log.Error(err)
	}
	return res
}

