package income

import (
	"database/sql"
	"fmt"
	"git.tor.ph/hiveon/pool/api/apierrors"
	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	secsInWeek = 7 * 24 * 3600
	secsInDay  = 24 * 3600
	wtfParam   = "2160 / 100000000.0"
)

// Repositorer represents common interface get miner profitability prediction and history
type IncomeRepositorer interface {
	GetIncome24h() (result float64, err error)
	GetIncome7d() (result float64, err error)
	GetIncomeResult() (result float64, err error)
	GetIncomeHistory() (*sql.Rows, error)
}

// Repository have methods to get miner income history and prediction
type IncomeRepository struct {
	client *gorm.DB
}

// NewRepository returns BlockRepository instance. DB MUST BE Mysqlnitize2
func NewIncomeRepository(db *gorm.DB) *IncomeRepository {
	return &IncomeRepository{db}
}

// GetIncome24h returns miner income for 24h
func (m *IncomeRepository) GetIncome24h() (result float64, err error) {
	sql :=
		fmt.Sprintf(`SELECT avg(expected_earning) * %s as income24h
		FROM expected_earning_result 
		WHERE end_time >= (UNIX_TIMESTAMP() -%d)`, wtfParam, secsInDay)

	if err := m.client.Raw(sql).Row().Scan(&result); apierrors.HandleError(err) {
		return 0, err
	}

	return result, nil
}

// GetIncome7d returns miner income for 7d
func (m *IncomeRepository) GetIncome7d() (result float64, err error) {
	sql :=
		fmt.Sprintf(`SELECT avg(expected_earning) * %s as income24h
		FROM expected_earning_result 
		WHERE end_time >= (UNIX_TIMESTAMP() -%d)`, wtfParam, secsInWeek)

	if err := m.client.Raw(sql).Row().Scan(&result); apierrors.HandleError(err) {
		return 0, err
	}

	return result, nil
}

func (m *IncomeRepository) GetIncomeResult() (result float64, err error) {
	sql :=
		fmt.Sprintf(`SELECT expected_earning * %s as income 
		FROM expected_earning_result 
		WHERE end_time >= UNIX_TIMESTAMP() limit 1`, wtfParam)

	if err := m.client.Raw(sql).Row().Scan(&result); apierrors.HandleError(err) {
		return 0, err
	}

	return result, nil
}

func (m *IncomeRepository) GetIncomeHistory() (*sql.Rows, error) {
	sql :=
		fmt.Sprintf(`SELECT start_time,expected_earning as income 
		FROM expected_earning_result 
		WHERE start_time >= (UNIX_TIMESTAMP() - %d)`, secsInDay)

	rows, err := m.client.Raw(sql).Rows()
	if apierrors.HandleError(err) {
		return nil, err
	}

	return rows, nil
}
