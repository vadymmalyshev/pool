package main

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"os"
	"os/signal"
	"syscall"

	"git.tor.ph/hiveon/pool/internal/billing"
	"github.com/sirupsen/logrus"
)

func main() {
	errs := make(chan error, 0)

	admDB, err := config.Config.Admin.DB.Connect()
	if err != nil {
		logrus.Panicf("failed to init Admin DB: %s", err)
	}

	calc := billing.NewBillingCalculator(admDB)
	calc.StartCalculation(errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Info("Billing terminated", <-errs)
}
