package users

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
)

type UserRepositorer interface {
	GetUserWallets(userID uint)            ([]models.Wallet, error)
	SaveUserWallet(user models.Wallet)     (models.Wallet, error)
	CreateCoinIfNotExists(coinName string) (*models.Coin, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositorer {
	return &UserRepository{db}
}

func (g *UserRepository) GetUserWallets(userID uint) ([]models.Wallet, error) {
	var userWallets []models.Wallet
	err := g.db.Where("user_id = ?", userID).Find(&userWallets).Error
	if apierrors.HandleError(err) {
		return []models.Wallet{}, err
	}
	return userWallets, nil
}

func (g *UserRepository) SaveUserWallet(user models.Wallet) (models.Wallet, error) {
	err := g.db.Save(&user).Error
	if apierrors.HandleError(err) {
		return models.Wallet{}, err
	}
	return user, nil
}

func (g *UserRepository) CreateCoinIfNotExists(coinName string) (*models.Coin, error) {
	var coin models.Coin
	err := g.db.FirstOrCreate(&coin, models.Coin{Name: coinName}).Error
	if apierrors.HandleError(err) {
		return nil, err
	}
	return &coin, nil
}
