package api

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/income"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	"github.com/gin-gonic/gin"
	"log"
)

type IncomeAPI struct {
	incomeService IncomeServicer
}

func NewIncomeAPI() *IncomeAPI {
	Sequelize2DB, err := mysql.Connect(config.Sequelize2DB)

	if err != nil {
		log.Panic("failed to init mysql Sequelize2DB db :", err.Error())
	}

	Sequelize3DB, err := mysql.Connect(config.Sequelize3DB)

	if err != nil {
		log.Panic("failed to init mysql Sequelize2DB db :", err.Error())
	}
	return &IncomeAPI{incomeService: NewIncomeService(Sequelize2DB, Sequelize3DB)}
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
