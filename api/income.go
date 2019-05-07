package api

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/internal/income"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IncomeAPI struct {
	incomeService income.IncomeServicer
}

func NewIncomeAPI(sql2DB *gorm.DB, sql3DB *gorm.DB) *IncomeAPI {
	return &IncomeAPI{incomeService: income.NewIncomeService(sql2DB, sql3DB)}
}

// Handle GET /api/pool/block/count24h
func (h *IncomeAPI) HandleGetBlockCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		bc, err := h.incomeService.GetBlockCount()
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, bc)
	}
}

// Handle GET /api/pool/incomeHistory
func (h *IncomeAPI) HandleGetIncomeHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		incHis, err := h.incomeService.GetIncomeHistory()
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, incHis)
	}
}
