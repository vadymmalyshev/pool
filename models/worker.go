package models

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"github.com/jinzhu/gorm"
)

const (
	tableNameWorker      = "workers"
	tableNameStatistic   = "billing_statistic"
)

type Worker struct {
	gorm.Model
	Name                string `gorm:"unique;not null"`
}

func (Worker) TableName() string {
	return tableNameWorker
}

func (Worker) AdminPath() string {
	return fmt.Sprintf("%s/%s", config.AdminPrefix, tableNameWorker)
}

type BillingWorkerStatistic struct {
	gorm.Model
	InvalidShares       float64
	StaleShares         float64
	ValidShares         float64
	ActivityPercentage  float64
	Worker              Worker `gorm:"foreignkey:WorkerID"`
	WorkerID            uint
	Wallet              Wallet `gorm:"foreignkey:WalletID"`
	WalletID            uint
}

func (BillingWorkerStatistic) TableName() string {
	return tableNameStatistic
}

func (BillingWorkerStatistic) AdminPath() string {
	return fmt.Sprintf("%s/%s", config.AdminPrefix, tableNameStatistic)
}
