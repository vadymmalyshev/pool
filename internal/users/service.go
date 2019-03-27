package users

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/models"
)

type UserServicer interface {
	GetUserWallet(userID uint) ([]Wallet, error)
	SaveUserWallet(userID uint, wallet string, coin string) (Wallet, error)
}

type userService struct {
	userRepository UserRepositorer
}

func NewUserService() UserServicer {
	return &userService{userRepository: NewUserRepository(config.Postgres)}
}

func (u *userService) GetUserWallet(userID uint) ([]Wallet, error) {
	return u.userRepository.GetUserWallets(userID)
}

func (u *userService) SaveUserWallet(userID uint, wallet string, coinName string) (Wallet, error) {
	coin, err := u.userRepository.CreateCoinIfNotExists(coinName)
	if err != nil {
		return Wallet{}, err
	}

	w := Wallet{Address: wallet, Coin: *coin}
	w, err = u.userRepository.SaveUserWallet(w)
	if err != nil {
		return Wallet{}, err
	}
	return w, nil
}
