package config

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var dbOnce sync.Once
var db, idpDB *gorm.DB

func initDatabases() {
	var err error
	db, err = Config.Admin.DB.Connect()
	if err != nil {
		logrus.Panicf("failed to init Admin db: %s", err)
	}

	idpDB, err = Config.IDP.DB.Connect()
	if err != nil {
		logrus.Panicf("failed to init IDP db: %s", err)
	}

	db.LogMode(true)
	idpDB.LogMode(true)
}

func InitPostgresDB() (*gorm.DB, *gorm.DB) {
	dbOnce.Do(func() {
		initDatabases()
	})

	return db, idpDB
}
