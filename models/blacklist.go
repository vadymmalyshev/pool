package models

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"time"
)

// Blacklist represents blacklist db model from Sequelize2DB
type Blacklist struct {
	ID          uint      `gorm:"column:id;primary_key"`
	Currency    string    `gorm:"column:currency"`
	BizNo       uint      `gorm:"column:biz_no"`
	MinerWallet string    `gorm:"column:miner_wallet"`
	OpType      string    `gorm:"column:op_type"`
	CreateDt    time.Time `gorm:"column:create_dt"`
}

const (
	tableNameBlacklist = "blacklist"
)

// TableName represent Blacklist table name. Used by Gorm
func (Blacklist) TableName() string {
	return tableNameBlacklist
}

func (Blacklist) AdminPath() string {
	return fmt.Sprintf("%s/%s", config.AdminPrefix, tableNameBlacklist)
}
