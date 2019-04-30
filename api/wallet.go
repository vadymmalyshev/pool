package api

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	walletsRepository "git.tor.ph/hiveon/pool/internal/wallets"
	"git.tor.ph/hiveon/pool/models"
	"github.com/gin-gonic/gin"
)

type EthAPI struct {
	service walletsRepository.WalletServicer
}

const (
	walletParam = "walletID"
	workerParam = "workerID"
)

// NewEthAPI return instance of ETH api service
func NewEthAPI() *EthAPI {
	return &EthAPI{service: walletsRepository.NewWalletService()}
}

// GetWalletFullData returns full mining history for 1d of specific wallet
func (api *EthAPI) GetWalletFullData() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(walletParam)
		wallInfo, err := api.service.GetWalletInfo(walletID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, wallInfo)
	}
}

// GetWalletsWorkerData returns mining history of specific worker of wallet
func (api *EthAPI) GetWalletsWorkerData() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(walletParam)
		workerID := c.Param(workerParam)
		info, err := api.service.GetWalletWorkerInfo(walletID, workerID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, info)
	}
}

// Handle POST /api/eth/wallet/addNew
func (api *EthAPI) AddNewWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var inputWallet models.Wallet
		err := c.BindJSON(&inputWallet)
		if apierrors.HandleError(err) {
			c.AbortWithStatusJSON(400, apierrors.NewApiErr(400, "Invalid JSON format"))
			return
		}
		w, err := api.service.AddWallet(inputWallet)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, w)
	}
}

//Handle DELETE /api/eth/:walletID/wallet/delete
func (api *EthAPI) DeleteWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(walletParam)
		err := api.service.DeleteWallet(walletID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.Status(201)
	}
}
