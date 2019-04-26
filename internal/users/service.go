package users

import (
	"git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
)

type UserServicer interface {
	GetUserWallet(userID uint) ([]models.Wallet, error)
	SaveUserWallet(userID uint, wallet string, coin string) (models.Wallet, error)
}

type userService struct {
	userRepository UserRepositorer
}

func NewUserService(admDB *gorm.DB) UserServicer {
	return &userService{userRepository: NewUserRepository(admDB)}
}

func (u *userService) GetUserWallet(userID uint) ([]models.Wallet, error) {
	return u.userRepository.GetUserWallets(userID)
}

func (u *userService) SaveUserWallet(userID uint, wallet string, coinName string) (models.Wallet, error) {
	coin, err := u.userRepository.CreateCoinIfNotExists(coinName)
	if err != nil {
		return models.Wallet{}, err
	}

	w := models.Wallet{Address: wallet, Coin: *coin}
	w, err = u.userRepository.SaveUserWallet(w)
	if err != nil {
		return models.Wallet{}, err
	}
	return w, nil
}
