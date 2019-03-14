package accounting

import (
	"database/sql"
	"fmt"

	"git.tor.ph/hiveon/pool/config"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type AccointingRepositorer interface {
	GetNormalBlocks24h() int
	GetUncleBlocks24h() int
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

// GetNormalBlocks24h returns a count of only normal blocks mined by the pool by last day
func (repo *AccountingRepository) GetNormalBlocks24h() int {
	sql := fmt.Sprintf(`
		SELECT count(id) as count 
		FROM blocks as b 
		WHERE is_uncle=0 AND 
		b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval %s)) * 1000`, config.PgOneDay)

	return repo.queryIntSingle(sql)
}

// GetUncleBlocks24h returns a count of only normal blocks mined by the pool by last day
func (repo *AccountingRepository) GetUncleBlocks24h() int {
	sql := fmt.Sprintf(`
		SELECT count(id) as count 
		FROM blocks as b 
		WHERE is_uncle=1 AND 
		b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval %s)) * 1000`, config.PgOneDay)

	return repo.queryIntSingle(sql)
}

func (repo *AccountingRepository) GetBill(walletId string) *sql.Rows {
	sql := fmt.Sprintf(`
		SELECT p.id, pd.paid, p.status, p.create_ts, p.tx_hash 
		FROM payment_details pd 
		INNER JOIN payments p ON p.id = pd.id 
		WHERE pd.miner_wallet = %s 
		ORDER BY pd.id desc 
		LIMIT 30`, walletId)

	rows, err := repo.db.Raw(sql).Rows()

	if err != nil {
		log.Error(err)
	}

	return rows
}

func (repo *AccountingRepository) GetBillInfo(walletId string) RepoBillInfo {
	var totalPaid float64
	err := repo.db.Raw(`
		SELECT sum(paid) as totalPaid 
		FROM payment_details 
		WHERE miner_wallet = ?`, walletId).Row().Scan(&totalPaid)

	if err != nil {
		log.Error(err)
	}

	var payment Payment
	var firstTime, balance string
	err1 := repo.db.Raw(`
		SELECT paid,payment_id 
		FROM payment_details 
		WHERE miner_wallet = ? 
		ORDER BY id LIMIT 1`, walletId).Row().Scan(&payment.firstPaid, &payment.paymentId)
	if err1 != nil {
		log.Error(err)
	}

	repo.db.Raw(`
		SELECT create_ts 
		FROM payments 
		WHERE id = ?`, payment.paymentId).Row().Scan(&firstTime)

	if err != nil {
		log.Error(err)
	}

	repo.db.Raw(`
		SELECT balance 
		FROM deposits 
		WHERE miner_wallet = ?`, walletId).Row().Scan(&balance)

	if err != nil {
		log.Error(err)
	}

	return RepoBillInfo{Balance: balance, FirstPaid: payment.firstPaid, FirstTime: firstTime, TotalPaid: totalPaid}
}

func (repo *AccountingRepository) GetBalance(walletId string) float64 {
	var res float64

	err := repo.db.Raw(`
		SELECT balance 
		FROM deposits 
		WHERE miner_wallet = ?`, walletId).Row().Scan(&res)
	if err != nil {
		log.Error(err)
	}
	return res
}
