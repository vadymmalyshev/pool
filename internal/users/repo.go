package users

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	. "git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
)

type UserRepositorer interface {
	GetUserWallets(userID uint) ([]Wallet, error)
	SaveUserWallet(user Wallet) (Wallet, error)
	CreateCoinIfNotExists(coinName string) (*Coin, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositorer {
	return &UserRepository{db}
}

func (g *UserRepository) GetUserWallets(userID uint) ([]Wallet, error) {
	var userWallets []Wallet
	err := g.db.Where("user_id = ?", userID).Find(&userWallets).Error
	if apierrors.HandleError(err) {
		return []Wallet{}, err
	}
	return userWallets, nil
}

func (g *UserRepository) SaveUserWallet(user Wallet) (Wallet, error) {
	err := g.db.Save(&user).Error
	if apierrors.HandleError(err) {
		return Wallet{}, err
	}
	return user, nil
}

func (g *UserRepository) CreateCoinIfNotExists(coinName string) (*Coin, error) {
	var coin Coin
	err := g.db.FirstOrCreate(&coin, Coin{Name: coinName}).Error
	if apierrors.HandleError(err) {
		return nil, err
	}
	return &coin, nil
}
