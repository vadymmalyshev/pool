package config

import (
	"flag"
	"os"

	"github.com/mitchellh/mapstructure"

	"git.tor.ph/hiveon/pool/internal/platform/hydra"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func YAMLUnmarshalOpt(c *mapstructure.DecoderConfig) {
	c.TagName = "yaml"
}

type AdminClient struct {
	ClientID     string `mapstructure:"admin.client_id"`
	ClientSecret string `mapstructure:"admin.client_secret"`
}

type API struct {
	Host string `mapstructure:"api.host"`
	Port int    `mapstructure:"api.port"`
}

type AdminDB struct {
	Host    string `mapstructure:"admin.db.host"`
	Port    int    `mapstructure:"admin.db.port"`
	SSLMode bool   `mapstructure:"admin.db.sslmode"`
	User    string `mapstructure:"admin.db.user"`
	Pass    string `mapstructure:"admin.db.password"`
	Name    string `mapstructure:"admin.db.name"`
	Log     string `mapstructure:"admin.db.log"`
}

// dbHost    = "admin.db.host"
// dbPort    = "admin.db.port"
// dbSSLMode = "admin.db.sslmode"
// dbUser    = "admin.db.user"
// dbPass    = "admin.db.password"
// dbName    = "admin.db.name"
// dbLog     = "admin.db.log"

// IDP represent IDP settings
type IDP struct {
	Host    string `mapstructure:"idp.db.host"`
	Port    int    `mapstructure:"idp.db.port"`
	SSLMode bool   `mapstructure:"idp.db.sslmode"`
	User    string `mapstructure:"idp.db.user"`
	Pass    string `mapstructure:"idp.db.password"`
	Name    string `mapstructure:"idp.db.name"`
	Log     string `mapstructure:"idp.db.log"`
}

// kafkaBrokers     = "kafka.brokers"
// kafkaCaLocation  = "kafka.ca_location"
// kafkaUsername    = "kafka.username"
// kafkaPass        = "kafka.password"
// kafkaTopics      = "kafka.topics"
// kafkaGroupId     = "kafka.group_id"
// kafkaRetention   = "kafka.retention"
// kafkaDbName      = "kafka.db_name"
// kafkaPrecision   = "kafka.precision"
// kafkaMiningPools = "kafka.mining_pools"

// AdminPrefix represents url prefix for admin panel
const AdminPrefix = "/admin"

var (
	AuthSignKey                                                                               string
	WorkerState, PoolZoom, ZoomConfigTime, ZoomConfigZoom, WorkerConfigTime, WorkerConfigZoom string
	HashrateCul, HashrateCulDivider                                                           string
	PgOneDay                                                                                  string
	MappingApi, WorkersAPI                                                                    string
	UseCasbin                                                                                 bool
	WorkerOfflineMin                                                                          int
	DefaultPercentage                                                                         float64

	DBConn, IDPDBConn string

	Hydra hydra.Config
)

var configPathFlag = *flag.String("c", "", "config file name from config directory")
var configPathEnv = os.Getenv("HIVEON_POOL_CONFIG")

func init() {

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if configPathFlag != "" {
		viper.SetConfigFile(configPathFlag)
	} else if configPathEnv != "" {
		viper.SetConfigFile(configPathEnv)
	} else {
		viper.AddConfigPath("$HOME/config")
		viper.AddConfigPath("./")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("$HIVEON_ADMIN_CONFIG_DIR/")

		viper.SetConfigName("config")
	}

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		logrus.Panicf("fatal error config file: %s", err)
	}

	if err := viper.Unmarshal(&Config, YAMLUnmarshalOpt); err != nil {
		logrus.Panicf("error while unmarshal viper config: %s", err)
		return
	}

	// Admin.Server = server.Config{
	// 	Host: viper.GetString("admin.host"),
	// 	Port: viper.GetInt("admin.port"),
	// }
	// if err := Admin.Server.Validate(); err != nil {
	// 	logrus.Panicf("Admin server configuration error: %s", err)
	// }

	// Admin.HydraClient = hydraclient.Config{
	// 	ClientID:     viper.GetString("admin.client_id"),
	// 	ClientSecret: viper.GetString("admin.client_secret"),
	// 	CallbackURL:  viper.GetString("admin.callback"),
	// }
	// if err := Admin.HydraClient.Validate(); err != nil {
	// 	logrus.Panicf("Admin server configuration error: %s", err)
	// }

	AuthSignKey = viper.GetString("auth.sign_key")

	if AuthSignKey == "" {
		panic("Token signing key is missing from configuration")
	}
	if len(AuthSignKey) < 32 {
		panic("Token signing key must be at least 32 characters")
	}
	WorkerOfflineMin = viper.GetInt("app.config.pool.workerOfflineMin")
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
	DefaultPercentage = viper.GetFloat64("WORKER_STAT_CONFIG.defaultPercentage")
	PgOneDay = viper.GetString("app.config.pool.pgOneDay")
	MappingApi = viper.GetString("pool.mapping_api")
	checkValueEmpty(MappingApi)
	WorkersAPI = viper.GetString("pool.workers_api")
	checkValueEmpty(WorkersAPI)
	checkValueEmpty(PgOneDay)
	UseCasbin = viper.GetBool("security.useCasbin")
	// influx
	AuthSignKey = viper.GetString("auth.sign_key")

	// if err := Redis.Validate(); err != nil {
	// 	logrus.Panicf("Redis configuration error: %s", err)
	// }

	// DB = database.Config{
	// 	Host:      viper.GetString(dbHost),
	// 	Port:      viper.GetInt(dbPort),
	// 	EnableSSL: viper.GetBool(dbSSLMode),
	// 	User:      viper.GetString(dbUser),
	// 	Pass:      viper.GetString(dbPass),
	// 	Name:      viper.GetString(dbName),
	// 	EnableLog: viper.GetBool(dbLog),
	// }

	// IDPDB = database.Config{
	// 	Host:      viper.GetString(idpdbHost),
	// 	Port:      viper.GetInt(idpdbPort),
	// 	EnableSSL: viper.GetBool(idpdbSSLMode),
	// 	User:      viper.GetString(idpdbUser),
	// 	Pass:      viper.GetString(idpdbPass),
	// 	Name:      viper.GetString(idpdbName),
	// 	EnableLog: viper.GetBool(idpdbLog),
	// }

	// Sequelize2DB = database.Config{
	// 	Host:      viper.GetString(sequelize2DBHost),
	// 	Port:      viper.GetInt(sequelize2DBPort),
	// 	EnableSSL: viper.GetBool(sequelize2DBSSLMode),
	// 	User:      viper.GetString(sequelize2DBUser),
	// 	Pass:      viper.GetString(sequelize2DBPass),
	// 	Name:      viper.GetString(sequelize2DBName),
	// }

	// Sequelize3DB = database.Config{
	// 	Host:      viper.GetString(sequelize3DBHost),
	// 	Port:      viper.GetInt(sequelize3DBPort),
	// 	EnableSSL: viper.GetBool(sequelize3DBSSLMode),
	// 	User:      viper.GetString(sequelize3DBUser),
	// 	Pass:      viper.GetString(sequelize3DBPass),
	// 	Name:      viper.GetString(sequelize3DBName),
	// }

	// InfluxDB = database.Config{
	// 	Host: viper.GetString(influxHost),
	// 	Port: viper.GetInt(influxPort),
	// 	User: viper.GetString(influxUser),
	// 	Pass: viper.GetString(influxPass),
	// 	Name: viper.GetString(influxName),
	// }

	// API = api.Config{
	// 	Host: viper.GetString(apiHost),
	// 	Port: viper.GetInt(apiPort),
	// }
}

func checkValueEmpty(val string) {
	if val == "" {
		panic(val + " missing from configuration")
	}
}
