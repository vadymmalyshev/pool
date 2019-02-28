package service

import (
	. "git.tor.ph/hiveon/pool/internal/api/repository"
	. "git.tor.ph/hiveon/pool/internal/api/response"
)
type UserService interface {
	GetUserWallet(userID uint) []UserWallets
	SaveUserWallet(userID uint, wallet string, coin string)
}

type userService struct {
	userRepository IUserRepository
}

func NewUserService() UserService{
	return &userService{userRepository:NewUserRepository()}
}

func (u *userService) GetUserWallet(userID uint) []UserWallets {
	return u.userRepository.GetUserWallets(userID)
}

func (u *userService) SaveUserWallet(userID uint, wallet string, coin string) {
	w:= UserWallets{UserID: userID, Wallet:wallet, Coin:coin}
	u.userRepository.SaveUserWallet(w)
	return
}
