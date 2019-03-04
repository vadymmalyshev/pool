package api

import (
	. "git.tor.ph/hiveon/pool/internal/api/service"
	"github.com/gin-gonic/gin"
)

type PoolAPI struct {
	poolService PoolService
}

func NewPoolAPI() *PoolAPI {
	return &PoolAPI{poolService:NewPoolService()}
}

func (h *PoolAPI) HandleGetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.poolService.GetIndex())
	}
}

func (h *PoolAPI) HandleGetIncomeHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.poolService.GetIncomeHistory())
	}
}


