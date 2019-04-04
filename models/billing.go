package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// WorkerFee - struct for work-fee table.
type WorkerFee struct {
	gorm.Model
	WorkerID   int     `json:"id" gorm:"column:worker_id;not null"`
	Date       string  `gorm:"column:billing_date;index"`
	WalletAddr string  `json:"wal" gorm:"column:wallet_address"`
	CoinName   string  `json:"coin" gorm:"column:coin"`
	Coin       Coin    `gorm:"foreignkey:CoinName"`
	Amount     float64 `json:"amount" gorm:"column:amount"`
	Shares     float64 `json:"shares" gorm:"column:shares"`
	Paid       bool    `json:"paid" gorm:"column:paid;index"`
}

// Charge - input struct for counting fee.
type Charge struct {
	Date       time.Time   `json:"date"`
	WorkersFee []WorkerFee `json:"workers"`
}

type Bill struct {
	WalletAdds string
	Coin       string
	Amount     float64
	Workers    []BillWorker
}

type BillWorker struct {
	ID     string
	Shares float64
}
