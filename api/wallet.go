package api

import (
	. "git.tor.ph/hiveon/pool/internal/api/service"
	"github.com/gin-gonic/gin"
)

// EthAPI container gin router handlers to get eth wallet data
type EthAPI struct {
	service WalletService
}

const (
	walletParam = "walletID"
	workerParam = "workerID"
)

// NewEthAPI return instance of ETH api service
func NewEthAPI() *EthAPI {
	return &EthAPI{walletService: NewWalletService()}
}

// GetWalletFullData returns full mining history for 1d of specific wallet
func (api *EthAPI) GetWalletFullData() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(walletParam)
		c.JSON(200, api.service.GetWalletInfo(walletID))
	}
}

// GetWalletsWorkerData returns mining history of specific worker of wallet
func (api *EthAPI) GetWalletsWorkerData() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(walletParam)
		workerID := c.Param(workerParam)
		c.JSON(200, api.service.GetWalletWorkerInfo(walletID, workerID))
	}
}
