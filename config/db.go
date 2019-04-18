package config

import (
	"sync"

	"github.com/jinzhu/gorm"
	// "github.com/sirupsen/logrus"
)

var dbOnce sync.Once
var db *gorm.DB

func initDB() {
	// config := {}NewDBConfig()
	// err := config.Validate()
	// if err != nil {
	// 	logrus.Panic("invalid database config: ", err.Error())
	// }

	// db, err = postgres.Connect(config)
	// if err != nil {
	// 	logrus.Panic("failed to initialize db: ", err.Error())
	// }
}

// GetDB returns an initialized DB instance.
func GetDB() *gorm.DB {
	dbOnce.Do(initDB)

	return db
}

// NewDBConfig returns a new Hiveon DB configuration struct
// func NewDBConfig() database.Config {
// 	return database.Config{
// 		Host:      "",
// 		Port:      "",
// 		EnableSSL: "",
// 		User:      "",
// 		Pass:      "",
// 		Name:      "",
// 		EnableLog: "",
// 	}
// }
