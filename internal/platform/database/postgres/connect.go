package postgres

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"

	// blank import is used here for simplicity
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbOnce sync.Once
var pointerDB *gorm.DB

// Connect returns initialized connection to db
func (db *DB) Connect() (*gorm.DB, error) {

	dbOnce.Do(func() {
		if err := db.Validate(); err != nil {
			logrus.Panic("failed to initialize db: ", err.Error())
			return
		}

		pointerDB, err := gorm.Open("postgres", db.Connection())
		if err != nil {
			logrus.Panic("failed to initialize db: ", err.Error())
			return
		}

		pointerDB.LogMode(db.Log)
	})

	return pointerDB, nil
}

// Connection represents connection string
func (db *DB) Connection() string {
	ssl := "disable"
	if db.SSLMode {
		ssl = "enable"
	}
	return fmt.Sprintf("host=%s port=%d sslmode=%s user=%s dbname=%s password=%s ",
		db.Host,
		db.Port,
		ssl,
		db.User,
		db.Name,
		db.Pass)
}
