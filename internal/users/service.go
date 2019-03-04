package repository

import (
	// "git.tor.ph/hiveon/pool/config"
	// . "git.tor.ph/hiveon/pool/internal/api/response"
	"log"

	// "git.tor.ph/hiveon/pool/internal/platform/database/postgres"
	// "github.com/jinzhu/gorm"
)

type IUserRepository interface {
	GetUserWallets(userID uint) []UserWallets
	SaveUserWallet(user UserWallets)
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

func (g *UserRepository) GetUserWallets(userID uint) []UserWallets {
	var userWallets []UserWallets
	g.client.Where("user_id = ?", userID).Find(&userWallets)
	return userWallets
}

func (g *UserRepository) SaveUserWallet(user UserWallets) {
	g.client.Save(&user)
	return
}
