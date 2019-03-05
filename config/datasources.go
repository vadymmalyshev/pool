package config

import (
	"git.tor.ph/hiveon/pool/internal/platform/database/influx"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/influxdata/influxdb1-client"
	red "github.com/go-redis/redis"
	"log"
	"strconv"
)

var (
	Seq2, Seq3, Postgres *gorm.DB
	Influx *client.Client
	Red    *redis.Client
	err error
)

func init() {
	Seq2, err = mysql.Connect(Sequelize2DB)

	if err != nil {
		log.Panic("failed to init mysql Sequelize2DB db :", err.Error())
	}

	Seq3, err = mysql.Connect(Sequelize3DB)
	if err != nil {
		log.Panic("failed to init mysql Sequelize3DB db :", err.Error())
	}

	Influx, err = influx.Connect(InfluxDB)
	if err != nil {
		log.Panic("failed to init influx:", err.Error())
	}

	Postgres, err = postgres.Connect(DB)
	if err != nil {
		log.Panic("failed to init postgres:", err.Error())
	}

	DBName, _ := strconv.Atoi(Redis.Name)
	Red := red.NewClient(&red.Options{
		Addr:     Redis.Host + ":" + strconv.Itoa(Redis.Port),
		Password: Redis.Pass,
		DB:       DBName,
	})
	_, err := Red.Ping().Result()
	if err != nil {
		log.Panic("failed to init redis:", err.Error())
	}
}
