package casbin

import (
	"github.com/jinzhu/gorm"
)

type CasbinRuleRepositorer interface {
	GetCasbinRule(ruleID string) *CasbinRule
	CreateCasbinRule(rule *CasbinRule) *CasbinRule
	UpdateCasbinRule(rule *CasbinRule) *CasbinRule
	DeleteCasbinRule(ruleID string)
	Migrate()
}

type CasRuleRepository struct {
	db *gorm.DB
}

func (g *CasRuleRepository) Migrate() {
	g.db.AutoMigrate(&CasbinRule{})
}

func NewCasRuleRepository(db *gorm.DB) CasbinRuleRepositorer {
	return &CasRuleRepository{db}
}

func (g *CasRuleRepository) GetCasbinRule(ruleID string) *CasbinRule {
	rule := new(CasbinRule)
	g.db.First(rule, ruleID)
	if rule.ID == 0 {
		return nil
	}
	return rule
}

func (g *CasRuleRepository) CreateCasbinRule(rule *CasbinRule) *CasbinRule {
	g.db.Create(rule)
	return rule
}

func (g *CasRuleRepository) UpdateCasbinRule(rule *CasbinRule) *CasbinRule {
	dbRule := new(CasbinRule)
	g.db.First(dbRule, rule.ID)
	if dbRule.ID == 0 {
		return nil
	}
	g.db.Omit("created_at", "deleted_at").Save(rule)
	return rule
}

func (g *CasRuleRepository) DeleteCasbinRule(ruleID string) {
	rule := new(CasbinRule)
	g.db.First(rule, ruleID)
	if rule.ID != 0 {
		g.db.Delete(rule)
	}
}
