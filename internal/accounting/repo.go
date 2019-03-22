package accounting

import (
	"database/sql"
	"fmt"
	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/config"
	"github.com/jinzhu/gorm"
)

type AccointingRepositorer interface {
	GetNormalBlocks24h() (int, error)
	GetUncleBlocks24h() (int, error)
	GetBillInfo(walletId string) (RepoBillInfo, error)
	GetBill(walletId string) (*sql.Rows, error)
	GetBalance(walletId string) (float64, error)
}

type AccountingRepository struct {
	db *gorm.DB
}

func NewAccountingRepository(db *gorm.DB) AccointingRepositorer {
	return &AccountingRepository{db: db}
}

func (repo *AccountingRepository) queryIntSingle(query string) (int, error) {
	var result int
	row := repo.db.Raw(query).Row()
	err := row.Scan(&result)
	if apierrors.HandleError(err) {
		return 0, err
	}
	return result, nil
}

// GetNormalBlocks24h returns a count of only normal blocks mined by the pool by last day
func (repo *AccountingRepository) GetNormalBlocks24h() (int, error) {
	sql := fmt.Sprintf(`
		SELECT count(id) as count 
		FROM blocks as b 
		WHERE is_uncle=0 AND 
		b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval %s)) * 1000`, config.PgOneDay)

	return repo.queryIntSingle(sql)
}

// GetUncleBlocks24h returns a count of only normal blocks mined by the pool by last day
func (repo *AccountingRepository) GetUncleBlocks24h() (int, error) {
	sql := fmt.Sprintf(`
		SELECT count(id) as count 
		FROM blocks as b 
		WHERE is_uncle=1 AND 
		b.block_ts > UNIX_TIMESTAMP(DATE_SUB(now(), interval %s)) * 1000`, config.PgOneDay)

	return repo.queryIntSingle(sql)
}

func (repo *AccountingRepository) GetBill(walletId string) (*sql.Rows, error) {
	sql := fmt.Sprintf(`
		SELECT p.id, pd.paid, p.status, p.create_ts, p.tx_hash 
		FROM payment_details pd 
		INNER JOIN payments p ON p.id = pd.id 
		WHERE pd.miner_wallet = %s 
		ORDER BY pd.id desc 
		LIMIT 30`, walletId)

	rows, err := repo.db.Raw(sql).Rows()

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (repo *AccountingRepository) GetBillInfo(walletId string) (RepoBillInfo, error) {
	var totalPaid float64
	// TotalPaid: if not paid yet, may be empty
	err := repo.db.Raw(`
		SELECT sum(paid) as totalPaid 
		FROM payment_details 
		WHERE miner_wallet = ?`, walletId).Row().Scan(&totalPaid)

	if apierrors.HandleError(err) {
		// Skip first err, if there no data in the rows
	}

	var payment Payment
	var firstTime, balance string
	// FirstPaid: payouts once a day, may be empty
	err = repo.db.Raw(`
		SELECT paid,payment_id 
		FROM payment_details 
		WHERE miner_wallet = ? 
		ORDER BY id LIMIT 1`, walletId).Row().Scan(&payment.firstPaid, &payment.paymentId)

	if apierrors.HandleError(err) {
		if err != sql.ErrNoRows {
			return RepoBillInfo{}, err
		}
	}

	// FirstTime: may be empty if there were no payments
	err = repo.db.Raw(`
		SELECT create_ts 
		FROM payments 
		WHERE id = ?`, payment.paymentId).Row().Scan(&firstTime)

	if apierrors.HandleError(err) {
		if err != sql.ErrNoRows {
			return RepoBillInfo{}, err
		}
	}

	// Balance: may be empty if a small mining time
	err = repo.db.Raw(`
		SELECT balance 
		FROM deposits 
		WHERE miner_wallet = ?`, walletId).Row().Scan(&balance)

	if apierrors.HandleError(err) {
		if err != sql.ErrNoRows {
			return RepoBillInfo{}, err
		}
	}
	repBillInf := RepoBillInfo{Balance: balance, FirstPaid: payment.firstPaid, FirstTime: firstTime, TotalPaid: totalPaid}

	if repBillInf.isEmpty() {
		return RepoBillInfo{}, apierrors.NewApiErr(400, "Bad Request")
	}

	return repBillInf, nil
}

func (repo *AccountingRepository) GetBalance(walletId string) (float64, error) {
	var res float64
	// Balance: may be empty if a small mining time
	err := repo.db.Raw(`
		SELECT balance 
		FROM deposits 
		WHERE miner_wallet = ?`, walletId).Row().Scan(&res)
	if apierrors.HandleError(err) {
		return 0, err
	}
	return res, nil
}
