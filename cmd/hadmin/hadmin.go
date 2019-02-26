package main

import (
	"os"

	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/models"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:   "hadmin",
	Short: "Hiveon Admin server",
	Run:   runServer,
}

var cmdAdmin = &cobra.Command{
	Use:   "admin",
	Short: "Add/remove admin rights to user",
	Run:   addAdmin,
}

var cmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate db",
	Run:   doMigrate,
}

func doMigrate(cmd *cobra.Command, args []string) {
	err := models.Migrate(config.GetDB())

	if err != nil {
		logrus.Panicf("something went wrong: %s", err)
	}
}

func addAdmin(cmd *cobra.Command, args []string) {
	logrus.Info("an admin created")
}

func runServer(cmd *cobra.Command, args []string) {
	logrus.Info("hAdmin server launched")

	// r := gin.New()
	// store, err := redis.NewStore(10, "tcp", config.Redis.Connection(), "", []byte(secret))

	// if err != nil {
	// 	logrus.Fatalf("can't connect to redis: %s", err.Error())
	// }

	// db := config.IDPDB()
	// defer db.Close()

}

const secret = "33446a9dcf9ea060a0a6532b166da32f304af0de"

func init() {
	cmdRoot.AddCommand(cmdMigrate)
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		logrus.Infof("can't run admin server: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
