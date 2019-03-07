package main

import (
	"fmt"
	"strconv"

	"git.tor.ph/hiveon/pool/api"
	"github.com/casbin/casbin"
	"github.com/casbin/redis-adapter"
	"os"
	"os/signal"
	"syscall"

	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	errs := make(chan error, 2)

	r := gin.Default()


	incomeHandler := api.NewIncomeAPI()
	ethAPI := api.NewEthAPI()
	minerHandler := api.NewMinerAPI()
	userHandler := api.NewUserAPI()

	r.GET("/api/pool/index", minerHandler.HandleGetIndex())
	r.GET("/api/pool/incomeHistory", incomeHandler.HandleGetIncomeHistory())
	r.GET("/api/pool/futureIncome", minerHandler.GetFutureIncome())
	r.GET("/api/pool/block/count24h", incomeHandler.HandleGetBlockCount())

	r.GET("/api/eth/:walletID", ethAPI.GetWalletFullData())
	r.GET("/api/eth/:walletID/worker/:workerID", ethAPI.GetWalletsWorkerData())
	r.GET("/api/eth/:walletID/billInfo", minerHandler.GetBillInfo())
	r.GET("/api/eth/:walletID/bill", minerHandler.GetBill())
	r.GET("/api/eth/:walletID/shares", minerHandler.GetShares())
	r.GET("/api/eth/:walletID/hashrate", minerHandler.GetHashrate())
	r.GET("/api/eth/:walletID/workers/counts", minerHandler.GetCountHistory())
	r.GET("/api/eth/:walletID/workers/list", minerHandler.GetMiner())

	r.GET("/api/private/statistic/worker/:workerID", minerHandler.GetWorkerStatistic())
	r.GET("/api/private/statistic/wallet/:walletID", minerHandler.GetWalletStatistic())
	r.GET("/api/private/statistic/workers", minerHandler.GetAllWorkersStatistic())
	r.GET("/api/private/statistic/mapping", minerHandler.GetWalletsWorkersMapping())

	r.GET("/api/private/wallet/:fid", userHandler.GetUserWallet())
	r.POST("/api/private/wallet", userHandler.PostUserWallet())

	initCasbinMiddleware(r)

	go func() {
		errs <- r.Run(fmt.Sprintf(":%d", config.API.Port))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Info("terminated", <-errs)
}

func initCasbinMiddleware(r *gin.Engine) {
	a_redis := redisadapter.NewAdapter("tcp", config.Redis.Host+":"+strconv.Itoa(config.Redis.Port))
	e := casbin.NewEnforcer("internal/casbin/authz_model.conf", a_redis)
	r.Use(NewAuthorizer(e))
}
