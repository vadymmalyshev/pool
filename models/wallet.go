package models

import (
	"fmt"

	"git.tor.ph/hiveon/pool/config"
	"github.com/jinzhu/gorm"
)

// Wallet represents wallet db model
type Wallet struct {
	gorm.Model
	Address             string `gorm:"not null"`
	Coin                Coin
	CoinID              uint `gorm:"index:idx_coin_id"`
	Percentage          float64 `gorm:"not null"`
}

const (
	tableNameWallet = "wallets"
	tableNameCoin   = "coins"
)

// TableName represent Wallet table name. Used by Gorm
func (Wallet) TableName() string {
	return tableNameWallet
}

func (Wallet) AdminPath() string {
	return fmt.Sprintf("%s/%s", config.AdminPrefix, tableNameWallet)
}

// Coin represents coin db model
type Coin struct {
	gorm.Model
	Name string `gorm:"not null"`
}

// TableName represent Coin table name. Used by Gorm
func (Coin) TableName() string {
	return tableNameCoin
}

func (Coin) AdminPath() string {
	return fmt.Sprintf("%s/%s", config.AdminPrefix, tableNameCoin)
}
