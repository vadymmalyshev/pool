package income

import (
	"database/sql"
	"fmt"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
)

const (
	secsInWeek = 7 * 24 * 3600
	secsInDay  = 24 * 3600
	wtfParam   = "2160 / 100000000.0"
)

// Repositorer represents common interface get miner profitability prediction and history
type IncomeRepositorer interface {
	GetIncome24h() (result float64)
	GetIncome7d() (result float64)
	GetIncomeResult() (result float64)
	GetIncomeHistory() *sql.Rows
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
func (m *IncomeRepository) GetIncome24h() (result float64) {
	sql :=
		fmt.Sprintf(`SELECT avg(expected_earning) * %s as income24h
		FROM expected_earning_result 
		WHERE end_time >= (UNIX_TIMESTAMP() -%d)`, wtfParam, secsInDay)

	m.client.Raw(sql).Row().Scan(&result)
	return result
}

// GetIncome7d returns miner income for 7d
func (m *IncomeRepository) GetIncome7d() (result float64) {
	sql :=
		fmt.Sprintf(`SELECT avg(expected_earning) * %s as income24h
		FROM expected_earning_result 
		WHERE end_time >= (UNIX_TIMESTAMP() -%d)`, wtfParam, secsInWeek)

	m.client.Raw(sql).Row().Scan(&result)
	return result
}

func (m *IncomeRepository) GetIncomeResult() (result float64) {
	sql :=
		fmt.Sprintf(`SELECT expected_earning * %s as income 
		FROM expected_earning_result 
		WHERE end_time >= UNIX_TIMESTAMP() limit 1`, wtfParam)

	m.client.Raw(sql).Row().Scan(&result)
	return result
}

func (m *IncomeRepository) GetIncomeHistory() *sql.Rows {
	sql :=
		fmt.Sprintf(`SELECT start_time,expected_earning as income 
		FROM expected_earning_result 
		WHERE start_time >= (UNIX_TIMESTAMP() - %d)`, secsInDay)

	rows, err := m.client.Raw(sql).Rows()
	if err != nil {
		log.Error(err)
	}

	return rows
}
