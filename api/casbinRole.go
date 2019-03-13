package api

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/casbin"
	"github.com/gin-gonic/gin"
	"go/types"
)

const paramCRID = "ruleID"

type CasbinRuleAPI struct {
	casRuleRepository CasbinRuleRepositorer
}

func NewCasbinRuleAPI() *CasbinRuleAPI {
	return &CasbinRuleAPI{NewCasRuleRepository(config.GetIDPDB())}
}

func (h *CasbinRuleAPI) MigrateRule() {
	h.casRuleRepository.Migrate()
}

// Handle GET /api/rule/get/:ruleID
func (h *CasbinRuleAPI) GetCasbinRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ruleId := c.Param(paramCRID)
		rule := h.casRuleRepository.GetCasbinRule(ruleId)
		if rule != nil {
			c.JSON(200, h.casRuleRepository.GetCasbinRule(ruleId))
		} else {
			c.JSON(200, types.Struct{})
		}
	}
}

// Handle POST /api/rule/create
func (h *CasbinRuleAPI) CreateCasbinRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule CasbinRule
		c.BindJSON(&rule)
		c.JSON(201, h.casRuleRepository.CreateCasbinRule(&rule))
	}
}

// Handle PUT /api/rule/update
func (h *CasbinRuleAPI) UpdateCasbinRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule CasbinRule
		c.BindJSON(&rule)
		res := h.casRuleRepository.UpdateCasbinRule(&rule)
		if res != nil {
			c.JSON(201, res)
		} else {
			c.JSON(201, types.Struct{})
		}
	}
}

// Handle DELETE /api/rule/delete/:ruleID
func (h *CasbinRuleAPI) DeleteCasbianRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ruleId := c.Param(paramCRID)
		h.casRuleRepository.DeleteCasbinRule(ruleId)
		c.Status(200)
	}
}
