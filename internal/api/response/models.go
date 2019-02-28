package response

import "github.com/jinzhu/gorm"

type PoolData struct {
	Code int `json:"code"`
	Data struct {
		Hashrate struct {
			Time        string  `json:"time"`
			ValidShares float64 `json:"validShares"`
			Hashrate    float64 `json:"hashrate"`
		} `json:"hashrate"`
		Miner struct {
			Time  string  `json:"time"`
			Count float64 `json:"count"`
		} `json:"miner"`
		Worker struct {
			Time  string  `json:"time"`
			Count float64 `json:"count"`
		} `json:"worker"`
	} `json:"data"`
}

// Success response
// swagger:response IncomeHistory
type IncomeHistory struct {
	Code int      `json:"code"`
	Data []Income `json:"data"`
}

type Income struct {
	Time   string `json:"time"`
	Income string `json:"income"`
}

//Success response
// swagger:response UserWallets
type UserWallets struct {
	gorm.Model
	Wallet string    `gorm:"not null;unique"`
	Coin   string    `gorm:"not null"`
	User   OAuthUser `gorm:"foreignkey:UserID"`
	UserID uint      `gorm:"not null"`
}

type OAuthUser struct {
	gorm.Model
	Username string `gorm:"not null"`
	Email string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Token string
	Challenge string
	Active bool `gorm:"not null"`
}