package config

import (
	"sync"

	"github.com/jinzhu/gorm"
)

var dbIdpOnce sync.Once
var idpDb *gorm.DB

func initIDPDB() {
	// config := NewIDPDBConfig()
	// err := config.Validate()
	// if err != nil {
	// 	logrus.Panic("invalid database config: ", err.Error())
	// }

	// idpDb, err = postgres.Connect(config)
	// if err != nil {
	// 	logrus.Panic("failed to initialize db: ", err.Error())
	// }
}

// GetIDPDB returns an initialized IDP DB instance.
// func GetIDPDB() *gorm.DB {
// 	// dbIdpOnce.Do(initIDPDB)
// 	// return idpDb
// }

// NewIDPDBConfig returns a new Admin DB configuration struct
// func NewIDPDBConfig() database.Config {
// 	return database.Config{
// 		Host:      viper.GetString(idpdbHost),
// 		Port:      viper.GetInt(idpdbPort),
// 		EnableSSL: viper.GetBool(idpdbSSLMode),
// 		User:      viper.GetString(idpdbUser),
// 		Pass:      viper.GetString(idpdbPass),
// 		Name:      viper.GetString(idpdbName),
// 		EnableLog: viper.GetBool(idpdbLog),
// 	}
// }
