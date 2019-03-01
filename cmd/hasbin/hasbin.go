package main

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/casbin"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errs := make(chan error, 0)

	sync, err := casbin.NewSynchronizer(config.DB, config.Redis)

	if err != nil {
		logrus.Panicf("can't start synchronizer: %s", err)
	}

	sync.Start(errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Info("Casbin sync terminated", <-errs)
}
