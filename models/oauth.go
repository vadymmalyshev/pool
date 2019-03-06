package models

import "github.com/jinzhu/gorm"

type OAuthUser struct {
	gorm.Model
	Username string `gorm:"not null"`
	Email string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Token string
	Challenge string
	Active bool `gorm:"not null"`
}
