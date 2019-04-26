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

// AdminPrefix represents url prefix for admin panel
const AdminPrefix = "/admin"

var (
	MappingApi, WorkersAPI string
	UseCasbin              bool

	DBConn, IDPDBConn string

	Hydra hydra.Config
)

const (
	SharesPerMHash = 4000000000 / 1200
	Devfee         = 0.03
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

	UseCasbin = viper.GetBool("security.useCasbin")

}

func checkValueEmpty(val string) {
	if val == "" {
		panic(val + " missing from configuration")
	}
}
