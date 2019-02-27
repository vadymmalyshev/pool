package config

import (
	"os"
	"path"
	"runtime"
	"strings"

	"git.tor.ph/hiveon/pool/internal/platform/api"
	"git.tor.ph/hiveon/pool/internal/platform/database"
	"git.tor.ph/hiveon/pool/internal/platform/hydra"
	hydraclient "git.tor.ph/hiveon/pool/internal/platform/hydra/client"
	"git.tor.ph/hiveon/pool/internal/platform/redis"
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
	influxHost = "influx.host"
	influxPort = "influx.port"
	influxUser = "influx.user"
	influxPass = "influx.password"
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
	AuthSignKey string
	Redis       redis.Config
	DB, IDPDB   database.Config

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

	Redis = redis.Config{
		Host: viper.GetString("redis.host"),
		Port: viper.GetInt("redis.port"),
		DB:   viper.GetString("redis.db"),
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

	API = api.Config{
		Host: viper.GetString(apiHost),
		Port: viper.GetInt(apiPort),
	}
}
