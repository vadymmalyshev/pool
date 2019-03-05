package api

import (
	. "git.tor.ph/hiveon/pool/internal/wallets"
	"github.com/gin-gonic/gin"
)

type WalletAPI struct {
	walletService WalletServicer
}

func NewWalletAPI() *WalletAPI {
	return &WalletAPI{walletService:NewWalletService()}
}

func (h *WalletAPI) HandleGetWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletId")
		c.JSON(200, h.walletService.GetWalletInfo(walletID))
	}
}

func (h *WalletAPI) HandleGetWalletWorker() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletId")
		workerID := c.Param("workerId")
		c.JSON(200, h.walletService.GetWalletWorkerInfo(walletID, workerID))
	}
}
