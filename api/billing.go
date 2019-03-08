package api

import (
	. "git.tor.ph/hiveon/pool/internal/billing"
	"github.com/gin-gonic/gin"
)

type BillingAPI struct {
	billingRepository 	BillingRepositorer
}

func NewBillingAPI() BillingAPI {
	return BillingAPI{billingRepository: NewBillingRepository()}
}

func (h *BillingAPI) HandleGetWalletEarning() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		date := c.Param(paramDate)
		c.JSON(200, h.billingRepository.GetWalletEarning(walletID, date))
	}
}


