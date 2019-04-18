package mysql

import (
	"fmt"

	// initialize mysql driver, must be there to make gorm works
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

// Connect returns initialized connection to db
func (db *DB) Connect() (*gorm.DB, error) {
	if err := db.Validate(); err != nil {
		return nil, err
	}

	conn, err := gorm.Open("mysql", db.Connection())
	if err != nil {
		return nil, err
	}

	conn.LogMode(db.Log)

	return conn, nil
}

// Connection returns db connection string
func (db *DB) Connection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		db.User,
		db.Pass,
		db.Host,
		db.Port,
		db.Name)
}
