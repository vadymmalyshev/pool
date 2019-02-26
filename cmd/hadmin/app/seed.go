package app

import (
	"fmt"
	"log"

	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/platform/database"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const seedTable = "last_update"

type idContainer struct {
	id int
}

func isTableExists(db *gorm.DB, schemaName string, tableName string) bool {
	var table idContainer

	sql := "SELECT 1 as id FROM information_schema.tables WHERE  table_catalog = '" + schemaName + "' AND table_name = '" + tableName + "'"
	db.Raw(sql).Scan(&table)

	return table.id > 0
}

// CreateAndSeedTables do create and fill tables
func CreateAndSeedTables() {
	db, err := database.Connect(config.DB)
	if err != nil {
		logrus.Panicf("can't connect to db: %s", err)
	}

	if !isTableExists(db, config.DB.Name, seedTable) {
		logrus.Infof("no need to seed db – table %s exist", seedTable)
		return
	}

	logrus.Info(fmt.Sprintf("Creating table %s ...", seedTable))

	db.Exec(fmt.Sprintf("CREATE TABLE %s ( id serial, tstamp timestamp DEFAULT now())", seedTable))
	db.Exec(fmt.Sprintf("INSERT INTO %s VALUES(1, now())", seedTable))

	logrus.Info(fmt.Sprintf("Creating stored procedure for %s ...", seedTable))

	db.Exec(`
		CREATE FUNCTION log_last_changes() 
		RETURNS trigger AS $BODY$ 
		BEGIN 
			UPDATE last_update 
			SET tstamp = now() 
			WHERE id = 1; 
			RETURN NEW; 
		END; $BODY$ language plpgsql;`)

	logrus.Info("Creating trigger for last_update ..")

	db.Exec(fmt.Sprintf("CREATE TRIGGER %s AFTER INSERT OR UPDATE OR DELETE ON casbin_rules FOR EACH STATEMENT EXECUTE PROCEDURE log_last_changes()", seedTable))

	// check user in ab_users
	userID := 0
	var user idContainer

	sql := "select count(*) as id from ab_users where username = ?"
	db.Raw(sql, "'admin'").Scan(&user)

	if user.id > 0 {
		userID = user.id
	} else { // create
		sql := "INSERT INTO ab_users VALUES(null,now(),now(),null,'admin','admin@hiveon.net','$2a$10$lmWdGp8ZJsFz5wJ9X8fi7uZ95XTC6zcx/trmd/TBuR3znx6.egrVC',null,null,true)"
		db.Exec(sql)
		userID = 999
	}

	// check user in casbin_rules
	var res idContainer
	sql = "select count(*) as id from casbin_rules where id = ?"
	db.Raw(sql, userID).Scan(&res)

	if !(res.id > 0) {
		log.Println("Added new user to casbin_rules ", userID)
		sql := "INSERT INTO casbin_rules VALUES(999, now(), now(), null, 'p', '999', '/*', '*')"
		db.Exec(sql)
	}
}
