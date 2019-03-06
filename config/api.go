package config

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type APIConfig struct {
	// Postgres
	HiveonDB *gorm.DB
	Redis    *redis.Conn
}
