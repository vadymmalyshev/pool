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
	handler := api.New()

	r.GET("/api/pool/index", handler.HandleGetIndex())
	r.GET("/api/pool/incomeHistory", handler.HandleGetIncomeHistory())

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

