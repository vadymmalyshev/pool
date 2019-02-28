package repository

import (
	"database/sql"

	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/platform/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type IBlockRepository interface {
	GetIncome24h() float64
	GetIncome7d() float64
	GetIncomeResult() float64
	GetIncomeHistory() *sql.Rows
}

type BlockRepository struct {
	client *gorm.DB
}

func GetBlockRepositoryClient() *gorm.DB {
	db, err := database.Connect(config.Sequelize3DB)
	defer db.Close()
	if err != nil {
		log.Panic("failed to init mysql db :", err.Error())
	}
	return db
}

func NewBlockRepository() IBlockRepository {
	return &BlockRepository{client: GetBlockRepositoryClient()}
}

func (m *BlockRepository) queryFloatSingle(query string) float64 {
	var result float64
	row := m.client.Raw(query)
	row.Scan(&result)
	return result
}

func (m *BlockRepository) GetIncome24h() float64 {
	sql :=
		`SELECT avg(expected_earning) * 2160/100000000.0 as income24h
		FROM expected_earning_result 
		WHERE end_time >= (UNIX_TIMESTAMP() -(24 * 3600))`

	return m.queryFloatSingle(sql)
}

func (m *BlockRepository) GetIncome7d() float64 {
	sql :=
		`SELECT avg(expected_earning) * 2160/100000000.0 as income24h
		FROM expected_earning_result 
		WHERE end_time >= (UNIX_TIMESTAMP() -(7 * 24 * 3600))`
	return m.queryFloatSingle(sql)
}

func (m *BlockRepository) GetIncomeResult() float64 {
	sql :=
		`SELECT expected_earning * 2160/100000000.0 as income 
	FROM expected_earning_result 
	WHERE end_time >= UNIX_TIMESTAMP() limit 1`
	return m.queryFloatSingle(sql)
}

func (m *BlockRepository) GetIncomeHistory() *sql.Rows {
	sql :=
		`SELECT start_time,expected_earning as income 
		FROM expected_earning_result 
		WHERE start_time >= (UNIX_TIMESTAMP() -(24 * 3600))`

	if rows, err := m.client.Raw(sql).Rows(); err != nil {
		log.Error(err)
	}
	return rows
}
