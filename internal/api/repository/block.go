package repository

import (
	"database/sql"
	"git.tor.ph/hiveon/pool/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"git.tor.ph/hiveon/pool/internal/platform/database"
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

func GetBlockRepositoryClient() *gorm.DB{
	db, err := database.ConnectMySQL(config.Sequelize3DB)
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
	sql := "SELECT avg(expected_earning) * 2160/100000000.0 as income24h" +
		" FROM expected_earning_result where end_time >= (UNIX_TIMESTAMP() -(24 * 3600))"
	return m.queryFloatSingle(sql)
}

func (m *BlockRepository) GetIncome7d() float64 {
	sql := "SELECT avg(expected_earning) * 2160/100000000.0 as income24h" +
		" FROM expected_earning_result where end_time >= (UNIX_TIMESTAMP() -(7 * 24 * 3600))"
	return m.queryFloatSingle(sql)
}

func (m *BlockRepository) GetIncomeResult() float64 {
	sql := "SELECT expected_earning * 2160/100000000.0 as income FROM expected_earning_result where end_time >= UNIX_TIMESTAMP() limit 1"
	return m.queryFloatSingle(sql)
}

func (m *BlockRepository) GetIncomeHistory() *sql.Rows {
	sql := "select start_time,expected_earning as income from expected_earning_result where" +
		" start_time >= (UNIX_TIMESTAMP() -(24 * 3600))"
	rows, err := m.client.Raw(sql).Rows()
	if err != nil {
		log.Error(err)
	}
	return rows
}
