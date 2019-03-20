package config

import (
	"git.tor.ph/hiveon/pool/internal/platform/database"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var idpDb *gorm.DB

func initIDPDB() {
	config := NewIDPDBConfig()
	err := config.Validate()
	if err != nil {
		logrus.Panic("invalid database config: ", err.Error())
	}

	idpDb, err = postgres.Connect(config)
	if err != nil {
		logrus.Panic("failed to initialize db: ", err.Error())
	}
}

// GetIDPDB returns an initialized IDP DB instance.
func GetIDPDB() *gorm.DB {
	dbOnce.Do(initIDPDB)

	return idpDb
}

// NewIDPDBConfig returns a new Admin DB configuration struct
func NewIDPDBConfig() database.Config {
	return database.Config{
		Host:      viper.GetString(idpdbHost),
		Port:      viper.GetInt(idpdbPort),
		EnableSSL: viper.GetBool(idpdbSSLMode),
		User:      viper.GetString(idpdbUser),
		Pass:      viper.GetString(idpdbPass),
		Name:      viper.GetString(idpdbName),
		EnableLog: viper.GetBool(idpdbLog),
	}
}
