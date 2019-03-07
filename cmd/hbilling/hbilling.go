package main

import (
	"fmt"
	. "git.tor.ph/hiveon/pool/internal/billing"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errs := make(chan error, 0)

	calc := NewBillingCalculator()
	calc.StartCalculation(errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Info("Billing terminated", <-errs)
}

