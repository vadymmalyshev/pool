package casbin

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	"github.com/jinzhu/gorm"
)

type CasbinRuleRepositorer interface {
	GetCasbinRule(ruleID string) (CasbinRule, error)
	CreateCasbinRule(rule CasbinRule) (CasbinRule, error)
	UpdateCasbinRule(rule CasbinRule) (CasbinRule, error)
	DeleteCasbinRule(ruleID string) error
}

type CasRuleRepository struct {
	db *gorm.DB
}

func NewCasRuleRepository(db *gorm.DB) CasbinRuleRepositorer {
	return &CasRuleRepository{db}
}

func (g *CasRuleRepository) GetCasbinRule(ruleID string) (CasbinRule, error) {
	rule := CasbinRule{}
	err := g.db.First(&rule, ruleID).Error
	if apierrors.HandleError(err) {
		return CasbinRule{}, err
	}
	return rule, nil
}

func (g *CasRuleRepository) CreateCasbinRule(rule CasbinRule) (CasbinRule, error) {
	err := g.db.Create(&rule).Error
	if apierrors.HandleError(err) {
		return CasbinRule{}, err
	}
	return rule, nil
}

func (g *CasRuleRepository) UpdateCasbinRule(rule CasbinRule) (CasbinRule, error) {
	dbRule := CasbinRule{}
	err := g.db.First(&dbRule, rule.ID).Error
	if apierrors.HandleError(err) {
		return CasbinRule{}, err
	}
	err = g.db.Omit("created_at", "deleted_at").Save(&rule).Error
	if apierrors.HandleError(err) {
		return CasbinRule{}, err
	}
	return rule, nil
}

func (g *CasRuleRepository) DeleteCasbinRule(ruleID string) error {
	rule := new(CasbinRule)
	err := g.db.First(rule, ruleID).Error
	if apierrors.HandleError(err) {
		return err
	}
	err = g.db.Delete(rule).Error
	if apierrors.HandleError(err) {
		return err
	}
	return nil
}
