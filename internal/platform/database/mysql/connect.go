package mysql

import (
	"fmt"

	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"

	"git.tor.ph/hiveon/pool/internal/platform/database"
)

// Connect returns initialized connection to db
func Connect(c database.Config) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", Connection(c))
	if err != nil {
		return nil, err
	}

	db.LogMode(c.EnableLog)

	return db, nil
}

// Connection returns db connection string
func Connection(c database.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name)
}
