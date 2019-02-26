package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // blank import is used here for simplicity
)

// Connect returns initialized connection to db
func Connect(c Config) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", c.Connection())
	if err != nil {
		return nil, err
	}

	db.LogMode(c.EnableLog)

	return db, nil
}
