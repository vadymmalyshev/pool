package main

import (
	"fmt"
	"git.tor.ph/hiveon/pool/internal/casbin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errs := make(chan error, 0)

	sync, err := casbin.NewSynchronizer()

	if err != nil {
		logrus.Panicf("can't start synchronizer: %s", err)
		errs <- err
	}

	sync.Start(errs)
	log.Info("Started syncronizer ")
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Info("terminated ", <-errs)
}
