package api

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	. "git.tor.ph/hiveon/pool/internal/billing"
	"github.com/gin-gonic/gin"
)

type BillingAPI struct {
	billingRepository BillingRepositorer
}

func NewBillingAPI() BillingAPI {
	return BillingAPI{billingRepository: NewBillingRepository()}
}

func (h *BillingAPI) HandleGetWalletEarning() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		date := c.Param(paramDate)
		wEarn, err := h.billingRepository.GetWalletEarning(walletID, date)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, wEarn)
	}
}
