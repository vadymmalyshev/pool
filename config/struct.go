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
	Admin        Admin        		`yaml:"admin"`
	SQL2         mysql.DB     		`yaml:"sequelize2"`
	SQL3         mysql.DB     		`yaml:"sequelize3"`
	InfluxDB     influx.DB    		`yaml:"influx"`
	Kafka        kafka.DB     		`yaml:"kafka"`
	Redis        redis.DB     		`yaml:"redis"`
	Hydra        hydra.Config 		`yaml:"hydra"`
	IDP          IDP          		`yaml:"idp"`
	Pool      	 Pool         		`yaml:"pool"`
	Scheduler 	 Scheduler    		`yaml:"scheduler"`
	API          API                `yaml:"api"`
}

// Scheduler represents settings for consumer scheduler
type Scheduler struct {
	Retention    string  		    `yaml:"eth_retention"`
	Measurement  string  			`yaml:"measurement"`
	EthAPI       string  			`yaml:"eth_api"`
	CnyAPI       string  			`yaml:"cny_api"`
}

type zoomAndTime struct {
	Period       string 			`yaml:"period"`
	Zoom         string 			`yaml:"zoom"`
}

type Hashrate  struct {
	Cul          string  			`yaml:"cul"`
	CulDivider   string  			`yaml:"cul_divider"`
}

type Blocks struct {
	Period       string  			`yaml:"period"`
}

type Workers struct {
	zoomAndTime          			`yaml:",squash"`
	OfflineAfter string  			`yaml:"offline_after"`
	State        string    			`yaml:"state"`
}

type Billing struct {
	DevFee       float64            `yaml:"dev_fee"`
}

type Admin struct {
	Server       serv.Config        `yaml:",squash"`
	Client       hydraclient.Config `yaml:",squash"`
	DB           postgres.DB        `yaml:"db"`
}

type Pool struct {
	WorkersAPI   string             `yaml:"workers_api"`
	MappingAPI   string      		`yaml:"mapping_api"`
	IdpAPI     	 string      		`yaml:"idp_api"`
	Shares       zoomAndTime 		`yaml:"shares"`
	Zoom       	 string      		`yaml:"zoom"`
	Hashrate     Hashrate    		`yaml:"hashrate"`
	Blocks       Blocks      		`yaml:"blocks"`
	Workers      Workers        	`yaml:"workers"`
	Billing      Billing         	`yaml:"billing"`
}

type API struct {
	Host       string      			`mapstructure:"api.host"`
	Port       int         			`mapstructure:"api.port"`
}

type IDP struct {
	Client     hydraclient.Config 	`yaml:",squash"`
	DB         postgres.DB        	`yaml:"db"`
}