package postgres

import (
	"fmt"

	"git.tor.ph/hiveon/pool/internal/platform/database"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // blank import is used here for simplicity
)

// Connect returns initialized connection to db
func Connect(c database.Config) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", Connection(c))
	if err != nil {
		return nil, err
	}

	db.LogMode(c.EnableLog)

	return db, nil
}

func Connection(c database.Config) string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "enable"
	}
	return fmt.Sprintf("host=%s port=%d sslmode=%s user=%s dbname=%s password=%s ",
		c.Host,
		c.Port,
		ssl,
		c.User,
		c.Name,
		c.Pass)
}
