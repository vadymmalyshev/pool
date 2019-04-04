package api

import (
	"time"

	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/internal/billing"
	"git.tor.ph/hiveon/pool/models"
	"github.com/gin-gonic/gin"
)

type BillingAPI struct {
	billingRepository billing.BillingRepositorer
	collector         *billing.Collector
}

func NewBillingAPI() BillingAPI {
	br := billing.NewBillingRepository()
	return BillingAPI{billingRepository: br, collector: billing.NewCollector(br)}
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

// HandleWorkersFee charges workers by Hiveon pool.
func (h *BillingAPI) HandleWorkersFee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var charge models.Charge
		err := c.BindJSON(&charge)
		if apierrors.HandleError(err) {
			c.AbortWithStatusJSON(400, apierrors.NewApiErr(400, "Invalid JSON format"))
			return
		}

		err = h.collector.ChargeWorkers(charge)
		if apierrors.AbortWithApiError(c, err) {
			return
		}

		c.JSON(201, charge.Date)
	}
}

// HandleWorkerBill updates workers bill.
func (h *BillingAPI) HandleWorkerBill() gin.HandlerFunc {
	return func(c *gin.Context) {
		date, err := time.Parse(time.RFC3339, c.Param(paramDate))
		if apierrors.HandleError(err) {
			c.AbortWithStatusJSON(400, apierrors.NewApiErr(400, "Invalid JSON format"))
			return
		}

		bills, err := h.collector.CalculateFees(date)
		if apierrors.AbortWithApiError(c, err) {
			return
		}

		c.JSON(201, bills)
	}
}
