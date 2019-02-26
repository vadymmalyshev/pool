package main

import (
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/casbin"
	"github.com/sirupsen/logrus"
)

func main() {
	errs := make(chan error, 0)

	sync, err := casbin.NewSynchronizer(config.DB, config.Redis)

	if err != nil {
		logrus.Panicf("can't start synchronizer: %s", err)
	}

	sync.Start(errs)
}
