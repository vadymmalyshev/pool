package api

import (
	. "git.tor.ph/hiveon/pool/internal/api/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PoolAPI struct {
	poolService PoolService
}

func New() *PoolAPI {
	return &PoolAPI{poolService:NewPoolService()}
}

func (h *PoolAPI) log() *logrus.Logger {
	return h.config.Log
}

func (h *PoolAPI) handleGetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.poolService.GetIndex())
	}
}

func (h *PoolAPI) handleGetIncomeHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.poolService.GetIncomeHistory())
	}
}


