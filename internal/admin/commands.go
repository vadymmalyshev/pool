package admin

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func AddAdmin(cmd *cobra.Command, args []string) {
	logrus.Info("an admin created")
}

func DoMigrate(cmd *cobra.Command, args []string) {
	adminDB, idpDB := config.InitPostgresDB()

	if err := models.Migrate(adminDB); err != nil {
		logrus.Panicf("something went wrong during migration to admin db: %s", err)
	}

	if err := models.MigrateWorkerFees(adminDB); err != nil {
		logrus.Panicf("something went wrong during migration to admin db: %s", err)
	}

	if err := models.MigrateIDP(idpDB); err != nil {
		logrus.Panicf("something went wrong during migration to idp db: %s", err)
	}
}

func RunServer(cmd *cobra.Command, args []string) {
	adminDB, idpDB := config.InitPostgresDB()
	seq2, err := config.Config.SQL2.Connect()

	if err != nil {
		logrus.Panicf("failed to init Seq2 db: %s", err)
	}

	logrus.Info("hAdmin server launched")

	mux := http.NewServeMux()
	InitAdmin(adminDB, idpDB, seq2, mux)

	r := gin.Default()
	r.Any("admin/*resources", gin.WrapH(mux))

	errs := make(chan error, 2)

	go func() {
		logrus.Infof("Hiveon Admin has started on https://%s", config.Config.Admin.Server.Addr())
		//errs <- r.RunTLS(config.Admin.Server.Addr(), config.Admin.Server.CertFile, config.Admin.Server.KeyFile)
		errs <- r.Run(config.Config.Admin.Server.Addr())
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Info("terminated", <-errs)
}