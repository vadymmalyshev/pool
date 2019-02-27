package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // blank import is used here for simplicity
	_ "github.com/go-sql-driver/mysql"
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

// Connect returns initialized connection to db
func ConnectMySQL(c Config) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", c.MySQLConnection())
	if err != nil {
		return nil, err
	}

	db.LogMode(c.EnableLog)

	return db, nil
}




