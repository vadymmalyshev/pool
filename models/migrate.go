package models

import (
	"fmt"
	"git.tor.ph/hiveon/pool/internal/casbin"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Migrate creates tables and relations by gorm models
func Migrate(db *gorm.DB) error {
	tables := []interface{}{
		&Wallet{},
		&Coin{},
		&Worker{},
		&BillingWorkerStatistic{},
		&BillingWorkerMoney{},
		&casbin.CasbinRule{},
	}

	var tableNames string
	for _, table := range tables {
		tableNames += fmt.Sprintf(" %s", db.NewScope(table).TableName())
	}

	logrus.WithFields(logrus.Fields{
		"table_names": strings.TrimSpace(tableNames),
	}).Info("migrating model tables")

	err := db.AutoMigrate(tables...).Error
	if err != nil {
		return err
	}

	err = AddForeignKey(db, &Coin{}, &Wallet{}, "CoinID")
	if err != nil {
		return err
	}

	return err
}

// AddForeignKeyAndReferencedKey supports add foreign key for gorm table models.
// Gorm doesn't support creating foreign keys when migrate tables.
func AddForeignKeyAndReferencedKey(db *gorm.DB, parentTable, childTable interface{}, foreignKeyField string, referencedField string) error {
	parentTableScope := db.NewScope(parentTable)
	childTableScope := db.NewScope(childTable)

	log := logrus.WithFields(logrus.Fields{
		"parent_table": strings.TrimSpace(parentTableScope.TableName()),
		"child_table":  strings.TrimSpace(childTableScope.TableName()),
	})

	f, ok := childTableScope.FieldByName(foreignKeyField)
	if !ok {
		return fmt.Errorf("field %q not found", foreignKeyField)
	}
	if !f.IsForeignKey {
		return fmt.Errorf("%q is not a foreign key field", foreignKeyField)
	}

	parentIdField := ""
	if referencedField == "" {
		parentIdField = parentTableScope.PrimaryKey()
	} else {
		f, ok := parentTableScope.FieldByName(referencedField)
		if !ok {
			return fmt.Errorf("field %q not found", referencedField)
		}
		parentIdField = f.DBName
	}
	references := fmt.Sprintf("%s(%s)", parentTableScope.TableName(), parentIdField)

	log.Infof("adding foreign key constraint: %s -> %s", f.DBName, references)
	return db.Model(childTable).AddForeignKey(f.DBName, references, "RESTRICT", "RESTRICT").Error
}

// AddForeignKey wraps AddForeignKeyAndReferencedKey with empty reference field
func AddForeignKey(db *gorm.DB, parentTable, childTable interface{}, foreignKeyField string) error {
	return AddForeignKeyAndReferencedKey(db, parentTable, childTable, foreignKeyField, "")
}
