package repository

import (
	"log"

	"git.tor.ph/hiveon/pool/config"
	."git.tor.ph/hiveon/pool/models"
	"git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	GetUserWallets(userID uint) []Wallet
	SaveUserWallet(user Wallet)
}

type UserRepository struct {
	client *gorm.DB
}

func GetUserRepositoryClient() *gorm.DB {
	db, err := postgres.Connect(config.DB)

	if err != nil {
		log.Panic("failed to init postgres db :", err.Error())
	}
	return db

}

func NewUserRepository() IUserRepository {
	return &UserRepository{client: GetUserRepositoryClient()}
}

func (g *UserRepository) GetUserWallets(userID uint) []Wallet {
	var userWallets []Wallet
	g.client.Where("user_id = ?", userID).Find(&userWallets)
	return userWallets
}

func (g *UserRepository) SaveUserWallet(user Wallet) {
	g.client.Save(&user)
	return
}
