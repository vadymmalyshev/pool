package api

import (
	. "git.tor.ph/hiveon/pool/internal/income"
	"github.com/gin-gonic/gin"
)

type IncomeAPI struct {
	incomeService IncomeServicer
}

func NewIncomeAPI() *IncomeAPI {
	return &IncomeAPI{incomeService: NewIncomeService()}
}

func (h *IncomeAPI) HandleGetBlockCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.incomeService.GetBlockCount())
	}
}

func (h *IncomeAPI) HandleGetIncomeHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.incomeService.GetIncomeHistory())
	}
}
