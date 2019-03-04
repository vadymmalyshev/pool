package users

import (
	"github.com/jinzhu/gorm"
	. "git.tor.ph/hiveon/pool/models"
)

type UserRepositorer interface {
	GetUserWallets(userID uint) []Wallet
	SaveUserWallet(user Wallet)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositorer {
	return &UserRepository{db}
}

func (g *UserRepository) GetUserWallets(userID uint) []Wallet {
	var userWallets []Wallet
	g.db.Where("user_id = ?", userID).Find(&userWallets)
	return userWallets
}

func (g *UserRepository) SaveUserWallet(user Wallet) {
	g.db.Save(&user)
	return
}
