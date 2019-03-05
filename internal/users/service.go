package users

import (
	. "git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
)

type UserServicer interface {
	GetUserWallet(userID uint) []Wallet
	SaveUserWallet(userID uint, wallet string, coin string)
}

type userService struct {
	userRepository UserRepositorer
}

func NewUserService(db *gorm.DB) UserServicer{
	return &userService{userRepository:NewUserRepository(db)}
}

func (u *userService) GetUserWallet(userID uint) []Wallet {
	return u.userRepository.GetUserWallets(userID)
}

func (u *userService) SaveUserWallet(userID uint, wallet string, coinName string) {
	coin, _ := u.userRepository.CreateCoinIfNotExists(coinName)

	w:= Wallet{UserID: userID, Address:wallet, Coin:*coin}
	u.userRepository.SaveUserWallet(w)
	return
}

