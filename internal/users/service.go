package users

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/models"
)

type UserServicer interface {
	GetUserWallet(userID uint) []Wallet
	SaveUserWallet(userID uint, wallet string, coin string)
}

type userService struct {
	userRepository UserRepositorer
}

func NewUserService() UserServicer{
	return &userService{userRepository:NewUserRepository(config.Postgres)}
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

