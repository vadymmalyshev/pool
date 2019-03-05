package users

import (
	"github.com/jinzhu/gorm"
	. "git.tor.ph/hiveon/pool/models"
)

type UserRepositorer interface {
	GetUserWallets(userID uint) []Wallet
	SaveUserWallet(user Wallet)
	CreateCoinIfNotExists(coinName string) (*Coin, error)
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

func (g *UserRepository) CreateCoinIfNotExists(coinName string) (*Coin, error) {
	var coin Coin
	res := g.db.FirstOrCreate(&coin, Coin{Name: coinName})
	if res.Error != nil {
		return nil, res.Error
	}
	return &coin, nil

}
