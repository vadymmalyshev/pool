package idp

import (
	"fmt"
	idp_users "git.tor.ph/hiveon/idp/models/users"
	"git.tor.ph/hiveon/pool/pkg/errors"
	"github.com/jinzhu/gorm"
)

type IDPRepositorer interface {
	GetUserID(email string) (uint, error)
}

type IDPRepository struct {
	db *gorm.DB
}

func NewIDPRepository(db *gorm.DB) IDPRepositorer {
	return &IDPRepository{db}
}

func (g *IDPRepository) GetUserID(email string) (uint, error) {
	var user idp_users.User
	notFoundByEmail := g.db.First(&user, "email = ?", email).RecordNotFound()

	if notFoundByEmail {
		fmt.Println("not found")
		return 0, errors.ErrUserNotFound
	}

	return user.ID, nil
}