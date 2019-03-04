package wallets

import (
	red "github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

// Config stores connections to databases
type Config struct {
	Redis *red.Conn
	DB    *gorm.DB
}

type service struct {
	Config
}

// NewWalletsService initialize service container by config
func NewWalletsService(c Config) *service {
	return &service{c}
}

// WalletService provides methods to get data by wallet id
type WalletService struct {
	s       *service
	Address string
}

func (s service) Wallet(walletID string) *WalletService {
	return &WalletService{&s, walletID}
}

// GetWorkersPulse return workers pulse map.
func (w WalletService) GetWorkersPulse() (map[string]WorkerPulse, error) {
	return GetWorkersPulse(*w.s.Redis, w.Address)
}
