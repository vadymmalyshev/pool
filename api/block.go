package api

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/income"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	"github.com/gin-gonic/gin"
	"log"
)

type BlockAPI struct {
	blockService IncomeServicer
}

func NewBlockAPI() *BlockAPI {
	db, err := mysql.Connect(config.Sequelize2DB)

	if err != nil {
		log.Panic("failed to init mysql Sequelize2DB db :", err.Error())
	}
	return &BlockAPI{blockService: NewIncomeService(db)}
}

func (h *BlockAPI) HandleGetBlockCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.blockService.GetBlockCount())
	}
}