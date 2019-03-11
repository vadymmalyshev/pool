package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.tor.ph/hiveon/idp/models/users"
	"git.tor.ph/hiveon/pool/config"
	internalAdmin "git.tor.ph/hiveon/pool/internal/admin"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	"git.tor.ph/hiveon/pool/models"

	"github.com/gin-gonic/gin"
	"github.com/qor/admin"
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
	db, err := postgres.Connect(config.DB)

	if err != nil {
		logrus.Panicf("failed to init db: %s", err)
	}

	err = models.Migrate(db)

	if err != nil {
		logrus.Panicf("something went wrong: %s", err)
	}
}

func addAdmin(cmd *cobra.Command, args []string) {
	logrus.Info("an admin created")
}

func runServer(cmd *cobra.Command, args []string) {
	db, err := postgres.Connect(config.DB)

	if err != nil {
		logrus.Panicf("failed to init hiveon db: %s", err)
	}

	idpdb, err := postgres.Connect(config.IDPDB)

	if err != nil {
		logrus.Panicf("failed to init idp db: %s", err)
	}

	seq2, err := mysql.Connect(config.Sequelize2DB)

	if err != nil {
		logrus.Panicf("failed to init seq2 db: %s", err)
	}

	logrus.Info("hAdmin server launched")

	admin := admin.New(&admin.AdminConfig{DB: idpdb})
	admin.GetRouter().Use(internalAdmin.SwitchDatabasesMiddleware(db, seq2))

	admin.AddResource(&models.Wallet{})
	admin.AddResource(&models.Coin{})
	admin.AddResource(&users.User{})
	admin.AddResource(&models.Blacklist{})

	mux := http.NewServeMux()
	admin.MountTo("/admin", mux)

	r := gin.Default()

	r.Any("admin/*resources", gin.WrapH(mux))

	errs := make(chan error, 2)

	go func() {
		logrus.Infof("Hiveon Admin has started on https://%s", config.Admin.Server.Addr())
		errs <- r.RunTLS(config.Admin.Server.Addr(), config.Admin.Server.CertFile, config.Admin.Server.KeyFile)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Info("terminated", <-errs)
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
