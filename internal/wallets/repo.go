package wallets

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
)

type WalletRepositorer interface {
	SaveWallet(models.Wallet) (models.Wallet, error)
	DeleteWallet(walletId string) error
}

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db}
}

func (r *WalletRepository) SaveWallet(wallet models.Wallet) (models.Wallet, error) {
	err := r.db.Create(&wallet).Error
	if apierrors.HandleError(err) {
		return models.Wallet{}, err
	}
	return wallet, nil
}

func (r *WalletRepository) DeleteWallet(wId string) error {
	var wallet models.Wallet
	err := r.db.First(&wallet, wId).Error
	if apierrors.HandleError(err) {
		return err
	}
	err = r.db.Delete(&wallet).Error
	if apierrors.HandleError(err) {
		return err
	}
	return nil
}
