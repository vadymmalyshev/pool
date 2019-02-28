package main

import (
	"fmt"
	"git.tor.ph/hiveon/pool/api"

	"os"
	"os/signal"
	"syscall"

	"git.tor.ph/hiveon/pool/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	errs := make(chan error, 2)

	r := gin.Default()
	poolHandler := api.NewPoolAPI()
	blockHandler := api.NewBlockAPI()
	walletHandler := api.NewWalletAPI()
	minerHandler := api.NewMinerAPI()
	userHandler := api.NewUserAPI()

	r.GET("/api/pool/index", poolHandler.HandleGetIndex())
	r.GET("/api/pool/incomeHistory", poolHandler.HandleGetIncomeHistory())

	r.GET("/api/block/count24h", blockHandler.HandleGetBlockCount())

	r.GET("/api/miner/ETH/{walletId}", walletHandler.HandleGetWallet())
	r.GET("/api/miner/ETH/{walletId}/{workerId}", walletHandler.HandleGetWalletWorker())

	r.GET("/api/miner/featureIncome", minerHandler.GetFutureIncome())
	r.GET("/api/miner/{walletId}/billInfo",minerHandler.GetBillInfo())
	r.GET("/api/miner/{walletId}/bill",minerHandler.GetBill())
	r.GET("/api/miner/{walletId}/shares",minerHandler.GetShares())
	r.GET("/api/miner/{walletId}/hashrate",minerHandler.GetHashrate())
	r.GET("/api/miner/{walletId}/workers/counts",minerHandler.GetCountHistory())
	r.GET("/api/miner/{walletId}/workers",minerHandler.GetMiner())

	r.GET("/api/private/statistic/worker/{workerId}",minerHandler.GetWorkerStatistic())
	r.GET("/api/private/statistic/wallet/{walletId}",minerHandler.GetWalletStatistic())
	r.GET("/api/private/statistic/workers",minerHandler.GetAllWorkersStatistic())
	r.GET("/api/private/statistic/mapping",minerHandler.GetWalletsWorkersMapping())

	r.GET("/api/private/{fid}", userHandler.GetUserWallet())
	r.POST("/api/private/wallet", userHandler.PostUserWallet())

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

