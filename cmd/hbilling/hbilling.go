package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.tor.ph/hiveon/pool/internal/billing"
	"github.com/sirupsen/logrus"
)

func main() {
	errs := make(chan error, 0)

	calc := billing.NewBillingCalculator()
	calc.StartCalculation(errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Info("Billing terminated", <-errs)
}
