package wallets

import (
	"errors"
	"fmt"
	"git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type WalletRepositorer interface {
	SaveWallet(*models.Wallet) (*models.Wallet, error)
	DeleteWallet(walletId string) error
}

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db}
}

func (r *WalletRepository) SaveWallet(wallet *models.Wallet) (*models.Wallet, error) {
	err := r.db.Create(wallet).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return wallet, nil
}

func (r *WalletRepository) DeleteWallet(wId string) error {
	var wallet models.Wallet
	if err := r.db.First(&wallet, wId).Error; err != nil {
		return err
	}
	if wallet.ID == 0 {
		return errors.New(fmt.Sprintf("invalid wallet ID:%s", wId))
	}
	if err := r.db.Delete(&wallet).Error; err != nil {
		return err
	}
	return nil
}
