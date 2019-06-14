package admin

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/casbin"
	"git.tor.ph/hiveon/pool/internal/idp"
	"git.tor.ph/hiveon/pool/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var casRuleRepository casbin.CasbinRuleRepositorer
var IDPRepository idp.IDPRepositorer
var adminDB, idpDB, seq2 *gorm.DB

func init() {
	var err error

	adminDB, idpDB = config.InitPostgresDB()
	seq2, err = config.Config.SQL2.Connect()

	if err != nil {
		logrus.Panicf("failed to init Seq2 db: %s", err)
	}
	casRuleRepository = casbin.NewCasRuleRepository(adminDB)
	IDPRepository = idp.NewIDPRepository(adminDB)
}

func AddAdmin(cmd *cobra.Command, args []string) {
	logrus.Info("Specify the command add/remove")
}

func AddAction(cmd *cobra.Command, args []string) {
	cr, err := getCasbinRule(args)

	if err != nil {
		logrus.Info(err)
		return
	}
	err = casRuleRepository.AddIfNotExistCasbinRule(cr)
	if err != nil {
		logrus.Info("Can't add admin for user ", args[0])
		return
	}
	logrus.Info("Admin added for user: ", args[0])
}

func RemoveAction(cmd *cobra.Command, args []string) {
	cr, err := getCasbinRule(args)

	if err != nil {
		logrus.Info(err)
		return
	}
	casRuleRepository.RemoveIfExistCasbinRule(cr)

	logrus.Info("Admin removed for user: ", args[0])
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

func createCasbinRule(id uint) casbin.CasbinRule{
	cr := casbin.CasbinRule{PType: "p", V0: strconv.FormatUint(uint64(id), 10), V1: "/*", V2: "*" }
	return cr
}

func getCasbinRule(args []string) (casbin.CasbinRule, error) {
	email := args[0]
	id, err := IDPRepository.GetUserID(email)

	if err != nil {
		logrus.Info("User's email is incorrect")
		return casbin.CasbinRule{}, err
	}
	cr := createCasbinRule(id)
	return cr, nil
}