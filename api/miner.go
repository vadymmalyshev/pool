package api

import (
	. "git.tor.ph/hiveon/pool/internal/api/service"
	"github.com/gin-gonic/gin"
)

type MinerAPI struct {
	minerService MinerService
}

func NewMinerAPI() *MinerAPI {
	return &MinerAPI{minerService:NewMinerService()}
}

func (h *MinerAPI) HandleGetFutureIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.minerService.GetFutureIncome())
	}
}