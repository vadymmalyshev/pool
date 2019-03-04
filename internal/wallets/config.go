package wallets

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

// Config contains list of data connections user by wallets service
type Config struct {
	// Postgres
	HiveonDB *gorm.DB
	Redis    *redis.Conn
	// paypemtns and deposits db
	AccountingDB *gorm.DB
}
