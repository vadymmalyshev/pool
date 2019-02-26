package models

import "github.com/jinzhu/gorm"

// WalletModel represents wallet db model
type WalletModel struct {
	gorm.Model
	Address string `gorm:"not null"`
	Coin    CoinModel
	CoinID  uint `gorm:"index:idx_coin_id"`
}

const (
	tableNameWallet = "wallets"
	tableNameCoin   = "coins"
)

// TableName represent WalletModel table name. Used by Gorm
func (WalletModel) TableName() string {
	return tableNameWallet
}

// CoinModel represents coin db model
type CoinModel struct {
	gorm.Model
	Name string `gorm:"not null"`
}

// TableName represent CoinModel table name. Used by Gorm
func (CoinModel) TableName() string {
	return tableNameCoin
}
