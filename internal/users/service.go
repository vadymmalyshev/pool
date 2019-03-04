package users

import (
	. "git.tor.ph/hiveon/pool/internal/api/repository"
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
	return &userService{userRepository:NewUserRepository()}
}

func (u *userService) GetUserWallet(userID uint) []Wallet {
	return u.userRepository.GetUserWallets(userID)
}

func (u *userService) SaveUserWallet(userID uint, wallet string, coin string) {
	//w:= Wallet{UserID: userID, Address:wallet, Coin:coin}
	//u.userRepository.SaveUserWallet(w)
	return
}
