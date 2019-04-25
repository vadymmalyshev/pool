package api

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/casbin"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

const paramCRID = "ruleID"

type CasbinRuleAPI struct {
	casRuleRepository casbin.CasbinRuleRepositorer
}

func NewCasbinRuleAPI() *CasbinRuleAPI {
	db, err := config.Config.IDP.DB.Connect()
	if err != nil {
		logrus.Panicf("failed to init db: %s", err)
	}
	return &CasbinRuleAPI{casbin.NewCasRuleRepository(db)}
}

// Handle GET /api/rule/get/:ruleID
func (h *CasbinRuleAPI) GetCasbinRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ruleId := c.Param(paramCRID)
		rule, err := h.casRuleRepository.GetCasbinRule(ruleId)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, rule)
	}
}

// Handle POST /api/rule/create
func (h *CasbinRuleAPI) CreateCasbinRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule casbin.CasbinRule
		err := c.BindJSON(&rule)
		if apierrors.HandleError(err) {
			c.AbortWithStatusJSON(400, apierrors.NewApiErr(400, "Invalid JSON format"))
			return
		}
		rule, err = h.casRuleRepository.CreateCasbinRule(rule)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(201, rule)
	}
}

// Handle PUT /api/rule/update
func (h *CasbinRuleAPI) UpdateCasbinRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule casbin.CasbinRule
		err := c.BindJSON(&rule)
		if apierrors.HandleError(err) {
			c.AbortWithStatusJSON(400, apierrors.NewApiErr(400, "Invalid JSON format"))
			return
		}
		res, err := h.casRuleRepository.UpdateCasbinRule(rule)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(201, res)
	}
}

// Handle DELETE /api/rule/delete/:ruleID
func (h *CasbinRuleAPI) DeleteCasbianRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ruleId := c.Param(paramCRID)
		err := h.casRuleRepository.DeleteCasbinRule(ruleId)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.Status(201)
	}
}
