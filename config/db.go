package config

import (
	"sync"

	"git.tor.ph/hiveon/pool/internal/platform/database"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var dbOnce sync.Once
var db *gorm.DB

func initDB() {
	config := NewDBConfig()
	err := config.Validate()
	if err != nil {
		logrus.Panic("invalid database config: ", err.Error())
	}

	db, err = postgres.Connect(config)
	if err != nil {
		logrus.Panic("failed to initialize db: ", err.Error())
	}
}

// GetDB returns an initialized DB instance.
func GetDB() *gorm.DB {
	dbOnce.Do(initDB)

	return db
}

// NewDBConfig returns a new Hiveon DB configuration struct
func NewDBConfig() database.Config {
	return database.Config{
		Host:      viper.GetString(dbHost),
		Port:      viper.GetInt(dbPort),
		EnableSSL: viper.GetBool(dbSSLMode),
		User:      viper.GetString(dbUser),
		Pass:      viper.GetString(dbPass),
		Name:      viper.GetString(dbName),
		EnableLog: viper.GetBool(dbLog),
	}
}
