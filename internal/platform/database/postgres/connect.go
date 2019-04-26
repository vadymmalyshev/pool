package postgres

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"

	// blank import is used here for simplicity
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbOnce sync.Once
var pointerDB *gorm.DB

// Connect returns initialized connection to db
func (db *DB) Connect() (*gorm.DB, error) {
	var (
		pointerDB *gorm.DB
		dbErr     error
	)
		var err error

		if err = db.Validate(); err != nil {
			dbErr = errors.Wrap(err, "failed to validate db config")
			return nil, err
		}

		if pointerDB, err = gorm.Open("postgres", db.Connection()); err != nil {
			dbErr = errors.Wrap(err, "failed to initialize db")
			return nil, err
		}

	if dbErr == nil {
		pointerDB.LogMode(db.Log)
	} else {
		dbOnce = *new(sync.Once)
	}

	return pointerDB, dbErr
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
