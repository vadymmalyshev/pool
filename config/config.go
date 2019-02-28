package config

import (
	"os"
	"path"
	"runtime"
	"strings"

	"git.tor.ph/hiveon/pool/internal/platform/api"
	"git.tor.ph/hiveon/pool/internal/platform/database"
	"git.tor.ph/hiveon/pool/internal/platform/hydra"
	"git.tor.ph/hiveon/pool/internal/platform/hydra/client"
	"git.tor.ph/hiveon/pool/internal/platform/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	apiHost = "api.host"
	apiPort = "api.port"
)

const (
	dbHost    = "admin.db.host"
	dbPort    = "admin.db.port"
	dbSSLMode = "admin.db.sslmode"
	dbUser    = "admin.db.user"
	dbPass    = "admin.db.password"
	dbName    = "admin.db.name"
	dbLog     = "admin.db.log"
)

const (
	idpdbHost    = "idp.db.host"
	idpdbPort    = "idp.db.port"
	idpdbSSLMode = "idp.db.sslmode"
	idpdbUser    = "idp.db.user"
	idpdbPass    = "idp.db.password"
	idpdbName    = "idp.db.name"
	idpdbLog     = "idp.db.log"
)

const (
	sequelize2DBHost    = "sequelize2.host"
	sequelize2DBPort    = "sequelize2.port"
	sequelize2DBSSLMode = "sequelize2.sslmode"
	sequelize2DBUser    = "sequelize2.user"
	sequelize2DBPass    = "sequelize2.password"
	sequelize2DBName    = "sequelize2.name"
)

const (
	sequelize3DBHost    = "sequelize3.host"
	sequelize3DBPort    = "sequelize3.port"
	sequelize3DBSSLMode = "sequelize3.sslmode"
	sequelize3DBUser    = "sequelize3.user"
	sequelize3DBPass    = "sequelize3.password"
	sequelize3DBName    = "sequelize3.name"
)

const (
	influxHost = "influx.host"
	influxPort = "influx.port"
	influxUser = "influx.user"
	influxPass = "influx.password"
	influxName = "influx.name"
)

const (
	appPort = "app.port"
	appHost = "app.host"

	hydraURL          = "hydra.url"
	hydraClientID     = "hydra.client_id"
	hydraClientSecret = "hydra.client_secret"
)

// AdminPrefix represents url prefix for admin panel
const AdminPrefix = "/admin"

type admin struct {
	Server      server.Config
	HydraClient hydraclient.Config
}

var (
	AuthSignKey                                                                               string
	WorkerState, PoolZoom, ZoomConfigTime, ZoomConfigZoom, WorkerConfigTime, WorkerConfigZoom string
	HashrateCul, HashrateCulDivider                                                           string
	PgOneDay                                                                                  string
	UseCasbin																			      bool
	Redis                                                                                     database.Config
	DB, IDPDB, Sequelize2DB, Sequelize3DB, InfluxDB                                           database.Config

	DBConn, IDPDBConn string

	Admin admin
	Hydra hydra.Config

	API api.Config
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	hiveonPoolDir := path.Join(path.Dir(filename), "..")

	os.Setenv("HIVEON_POOL", hiveonPoolDir)

	viper.AddConfigPath("$HOME/config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HIVEON_ADMIN_CONFIG_DIR/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetConfigName("config")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		logrus.Panicf("Fatal error config file: %s", err)
	}

	Admin.Server = server.Config{
		Host:     viper.GetString("admin.host"),
		Port:     viper.GetInt("admin.port"),
		CertFile: strings.Replace(viper.GetString("admin.certs.pem"), "$HIVEON_POOL", hiveonPoolDir, -1),
		KeyFile:  strings.Replace(viper.GetString("admin.certs.key"), "$HIVEON_POOL", hiveonPoolDir, -1),
	}
	if err := Admin.Server.Validate(); err != nil {
		logrus.Panicf("Admin server configuration error: %s", err)
	}

	Admin.HydraClient = hydraclient.Config{
		ClientID:     viper.GetString("admin.client_id"),
		ClientSecret: viper.GetString("admin.client_secret"),
		CallbackURL:  viper.GetString("admin.callback"),
	}
	if err := Admin.HydraClient.Validate(); err != nil {
		logrus.Panicf("Admin server configuration error: %s", err)
	}

	AuthSignKey = viper.GetString("auth.sign_key")

	if AuthSignKey == "" {
		panic("Token signing key is missing from configuration")
	}
	if len(AuthSignKey) < 32 {
		panic("Token signing key must be at least 32 characters")
	}

	WorkerState = viper.GetString("app.config.pool.workerState")
	PoolZoom = viper.GetString("app.config.pool.poolZoom")
	ZoomConfigTime = viper.GetString("ZOOM_CONFIG.d.time")
	ZoomConfigZoom = viper.GetString("ZOOM_CONFIG.d.zoom")
	WorkerConfigTime = viper.GetString("WORKER_STAT_CONFIG.d.time")
	WorkerConfigZoom = viper.GetString("WORKER_STAT_CONFIG.d.zoom")

	HashrateCul = viper.GetString("app.config.pool.hashrate.hashrateCul")
	checkValueEmpty(HashrateCul)
	HashrateCulDivider = viper.GetString("app.config.pool.hashrate.hashrateCulDivider")
	checkValueEmpty(HashrateCulDivider)
	PgOneDay = viper.GetString("app.config.pool.pgOneDay")
	checkValueEmpty(PgOneDay)
	UseCasbin = viper.GetBool("security.useCasbin")
	// influx
	AuthSignKey = viper.GetString("auth.sign_key")

	Redis = database.Config{
		Host: viper.GetString("redis.host"),
		Port: viper.GetInt("redis.port"),
		Name:   viper.GetString("redis.db"),
		Pass: viper.GetString("redis.password"),
	}
	if err := Redis.Validate(); err != nil {
		logrus.Panicf("Redis configuration error: %s", err)
	}

	DB = database.Config{
		Host:      viper.GetString(dbHost),
		Port:      viper.GetInt(dbPort),
		EnableSSL: viper.GetBool(dbSSLMode),
		User:      viper.GetString(dbUser),
		Pass:      viper.GetString(dbPass),
		Name:      viper.GetString(dbName),
		EnableLog: viper.GetBool(dbLog),
	}

	IDPDB = database.Config{
		Host:      viper.GetString(idpdbHost),
		Port:      viper.GetInt(idpdbPort),
		EnableSSL: viper.GetBool(idpdbSSLMode),
		User:      viper.GetString(idpdbUser),
		Pass:      viper.GetString(idpdbPass),
		Name:      viper.GetString(idpdbName),
		EnableLog: viper.GetBool(idpdbLog),
	}

	Sequelize2DB = database.Config{
		Host:      viper.GetString(sequelize2DBHost),
		Port:      viper.GetInt(sequelize2DBPort),
		EnableSSL: viper.GetBool(sequelize2DBSSLMode),
		User:      viper.GetString(sequelize2DBUser),
		Pass:      viper.GetString(sequelize2DBPass),
		Name:      viper.GetString(sequelize2DBName),
	}

	Sequelize3DB = database.Config{
		Host:      viper.GetString(sequelize3DBHost),
		Port:      viper.GetInt(sequelize3DBPort),
		EnableSSL: viper.GetBool(sequelize3DBSSLMode),
		User:      viper.GetString(sequelize3DBUser),
		Pass:      viper.GetString(sequelize3DBPass),
		Name:      viper.GetString(sequelize3DBName),
	}

	InfluxDB = database.Config{
		Host: viper.GetString(influxHost),
		Port: viper.GetInt(influxPort),
		User: viper.GetString(influxUser),
		Pass: viper.GetString(influxPass),
		Name: viper.GetString(influxName),
	}

	API = api.Config{
		Host: viper.GetString(apiHost),
		Port: viper.GetInt(apiPort),
	}
}

func checkValueEmpty(val string) {
	if val == "" {
		panic(val + " missing from configuration")
	}
}
