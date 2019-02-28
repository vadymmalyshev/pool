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

	r.GET("/api/pool/index", poolHandler.HandleGetIndex())
	r.GET("/api/pool/incomeHistory", poolHandler.HandleGetIncomeHistory())

	r.GET("/api/block/count24h", blockHandler.HandleGetBlockCount())

	r.GET("/api/miner/ETH/{walletId}", walletHandler.HandleGetWallet())
	r.GET("/api/miner/ETH/{walletId}/{workerId}", walletHandler.HandleGetWalletWorker())

	r.GET("/api/miner/featureIncome", minerHandler.HandleGetFutureIncome())

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

