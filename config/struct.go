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

type zoomAndTime struct {
	Period string `yaml:"period"`
	Zoom   string `yaml:"zoom"`
}

type common struct {
	Admin struct {
		Server serv.Config        `yaml:",squash"`
		Client hydraclient.Config `yaml:",squash"`
		DB     postgres.DB        `yaml:"db"`
	} `yaml:"admin"`
	SQL2     mysql.DB     `yaml:"sequelize2"`
	SQL3     mysql.DB     `yaml:"sequelize3"`
	InfluxDB influx.DB    `yaml:"influx"`
	Kafka    kafka.DB     `yaml:"kafka"`
	Redis    redis.DB     `yaml:"redis"`
	Hydra    hydra.Config `yaml:"hydra"`
	IDP      struct {
		Client hydraclient.Config `yaml:",squash"`
		DB     postgres.DB
	} `yaml:"idp"`
	Pool struct {
		WorkersAPI string      `yaml:"workers_api"`
		MappingAPI string      `yaml:"mapping_api"`
		IdpAPI     string      `yaml:"idp_api"`
		Shares     zoomAndTime `yaml:"shares"`
		Zoom       string      `yaml:"zoom"`
		Blocks     struct {
			Period string `yaml:"period"`
		} `yaml:"blocks"`
		Workers struct {
			zoomAndTime  `yaml:",squash"`
			OfflineAfter string `yaml:"offline_after"`
			State        string `yaml:"state"`
		} `yaml:"shares"`
	} `yaml:"pool"`
	Scheduler  Scheduler `yaml:"scheduler"`
}

// Scheduler represents settings for consumer scheduler
type Scheduler struct {
	Retention   string `yaml:"eth_retention"`
	Measurement string `yaml:"measurement"`
	EthAPI      string `yaml:"eth_api"`
	CnyAPI      string `yaml:"cny_api"`
}
