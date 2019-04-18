package config

import (
	"git.tor.ph/hiveon/pool/internal/platform/database/influx"
	"git.tor.ph/hiveon/pool/internal/platform/database/kafka"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	"git.tor.ph/hiveon/pool/internal/platform/database/redis"
	"git.tor.ph/hiveon/pool/internal/platform/hydra"
	hydraclient "git.tor.ph/hiveon/pool/internal/platform/hydra/client"
	serv "git.tor.ph/hiveon/pool/internal/platform/server"
)

// Config represent all of the settings
var Config common

type common struct {
	Admin struct {
		serv.Config `yaml:",inline"`
		Client      hydraclient.Config `yaml:",inline"`
		DB          postgres.DB        `yaml:"db"`
	} `yaml:"admin"`
	SQL2     mysql.DB     `yaml:"sequelize2"`
	SQL3     mysql.DB     `yaml:"sequelize3"`
	InfluxDB influx.DB    `yaml:"influx"`
	Kafka    kafka.DB     `yaml:"kafka"`
	Redis    redis.DB     `yaml:"redis"`
	Hydra    hydra.Config `yaml:"hydra"`
	IDP      struct {
		Client hydraclient.Config
		DB     postgres.DB
	} `yaml:"idp"`
	WorkersAPI string `yaml:"pool.workers_api"`
	MappingAPI string `yaml:"pool.mapping_api"`
	IdpAPI     string `yaml:"pool.idp_api"`
}

// Scheduler represents settings for consumer scheduler
type Scheduler struct {
	Retention   string `yaml:"eth_retention"`
	Measurement string `yaml:"measurement"`
	EthAPI      string `yaml:"eth_api"`
	CnyAPI      string `yaml:"cny_api"`
}
