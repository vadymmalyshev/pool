package api

import (
	. "git.tor.ph/hiveon/pool/internal/api/service"
	"github.com/gin-gonic/gin"
)

type BlockAPI struct {
	blockService BlockService
}

func NewBlockAPI() *BlockAPI {
	return &BlockAPI{blockService:NewBlockService()}
}

func (h *BlockAPI) HandleGetBlockCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.blockService.GetBlockCount())
	}
}