package repository

import (
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepository{db}
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
