package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/casbin/casbin"
	"github.com/casbin/redis-adapter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"git.tor.ph/hiveon/pool/api"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/api/middleware"
)

func main() {
	errs := make(chan error, 2)

	r := gin.Default()

	incomeHandler := api.NewIncomeAPI()
	walletHandler := api.NewEthAPI()
	minerHandler := api.NewMinerAPI()
	userHandler := api.NewUserAPI()
	billingHandler := api.NewBillingAPI()

	casbinRuleHandler := api.NewCasbinRuleAPI()

	r.GET("/api/pool/index", minerHandler.HandleGetIndex())
	r.GET("/api/pool/incomeHistory", incomeHandler.HandleGetIncomeHistory())
	r.GET("/api/pool/futureIncome", minerHandler.GetFutureIncome())
	r.GET("/api/pool/block/count24h", incomeHandler.HandleGetBlockCount())

	r.GET("/api/eth/:walletID", walletHandler.GetWalletFullData())
	r.GET("/api/eth/:walletID/worker/:workerID", walletHandler.GetWalletsWorkerData())
	r.GET("/api/eth/:walletID/billInfo", minerHandler.GetBillInfo())
	r.GET("/api/eth/:walletID/bill", minerHandler.GetBill())
	r.GET("/api/eth/:walletID/shares", minerHandler.GetShares())
	r.GET("/api/eth/:walletID/hashrate", minerHandler.GetHashrate())
	r.GET("/api/eth/:walletID/workers/counts", minerHandler.GetCountHistory())
	r.GET("/api/eth/:walletID/workers/list", minerHandler.GetMiner())
	r.GET("/api/eth/:walletID/date/:date/earning", billingHandler.HandleGetWalletEarning())
	r.POST("/api/eth/wallet/addNew", walletHandler.AddNewWallet())
	r.DELETE("/api/eth/:walletID/wallet/delete", walletHandler.DeleteWallet())

	r.GET("/api/private/statistic/worker/:workerID", minerHandler.GetWorkerStatistic())
	r.GET("/api/private/statistic/wallet/:walletID", minerHandler.GetWalletStatistic())
	r.GET("/api/private/statistic/workers", minerHandler.GetAllWorkersStatistic())
	r.GET("/api/private/statistic/mapping", minerHandler.GetWalletsWorkersMapping())

	r.GET("/api/private/wallet/:fid", userHandler.GetUserWallet())
	r.POST("/api/private/wallet", userHandler.PostUserWallet())

	r.GET("/api/rule/get/:ruleID", casbinRuleHandler.GetCasbinRule())
	r.POST("/api/rule/create", casbinRuleHandler.CreateCasbinRule())
	r.PUT("/api/rule/update", casbinRuleHandler.UpdateCasbinRule())
	r.DELETE("/api/rule/delete/:ruleID", casbinRuleHandler.DeleteCasbianRule())

	r.POST("/api/charge", billingHandler.HandleWorkersFee())
	r.POST("/api/bill/:date", billingHandler.HandleWorkerBill())

	initCasbinMiddleware(r)

	go func() {
		errs <- r.Run(fmt.Sprintf(":%d", config.Config.API.Port))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Info("terminated", <-errs)
}

func initCasbinMiddleware(r *gin.Engine) {
	a_redis := redisadapter.NewAdapter("tcp", config.Config.Redis.Host+":"+strconv.Itoa(config.Config.Redis.Port))
	e := casbin.NewEnforcer("internal/casbin/authz_model.conf", a_redis)
	r.Use(middleware.NewAuthorizer(e))
}
