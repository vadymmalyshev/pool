package api

import (
	. "git.tor.ph/hiveon/pool/internal/wallets"
	"git.tor.ph/hiveon/pool/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type EthAPI struct {
	service WalletServicer
}

const (
	walletParam = "walletID"
	workerParam = "workerID"
)

// NewEthAPI return instance of ETH api service
func NewEthAPI() *EthAPI {
	return &EthAPI{service: NewWalletService()}
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

// Handle POST /api/eth/wallet/addNew
func (api *EthAPI) AddNewWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var inputWallet models.Wallet
		err := c.BindJSON(&inputWallet)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatusJSON(400, "Invalid JSON format")
			return
		}
		w, err := api.service.AddWallet(&inputWallet)
		if err != nil {
			c.AbortWithStatusJSON(400, "Api has error")
			return
		}
		c.JSON(200, w)
	}
}

//Handle DELETE /api/eth/:walletID/wallet/delete
func (api *EthAPI) DeleteWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(walletParam)
		if err := api.service.DeleteWallet(walletID); err != nil {
			c.AbortWithStatusJSON(400, "Api has error")
			return
		}
		c.Status(201)
	}
}
