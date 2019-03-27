package casbin

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/jinzhu/gorm"
)

const adminAccess = `
[request definition]
r = sub, obj, act

[policy definition]
p = sub, obj, act

[role definition]
g = _,_

[policy effect]
e = some(where (p.eft = allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && (r.act == p.act || p.act == "*")
`

// IsHaveAdminAccess returns true if user have access to admin
func IsHaveAdminAccess(a gormadapter.Adapter, uid string) bool {
	adminAccessEnforcer := casbin.NewEnforcer(a, casbin.NewModel(adminAccess))
	return adminAccessEnforcer.Enforce(uid, "admin", "access")
}

// CasbinRule represents casbin_rule db model
type CasbinRule struct {
	gorm.Model
	PType string `gorm:"type:varchar(100)" json:"p_type"`
	V0    string `gorm:"type:varchar(100)" json:"v0"`
	V1    string `gorm:"type:varchar(100)" json:"v1"`
	V2    string `gorm:"type:varchar(100)" json:"v2"`
	V3    string `gorm:"type:varchar(100)" json:"v3"`
	V4    string `gorm:"type:varchar(100)" json:"v4"`
	V5    string `gorm:"type:varchar(100)" json:"v5"`
}

const (
	tableNameCasbinRule = "casbin_rules"
)

// TableName represent CasbinRule table name. Used by Gorm
func (CasbinRule) TableName() string {
	return tableNameCasbinRule
}

func (CasbinRule) AdminPath() string {
	return fmt.Sprintf("%s/%s", config.AdminPrefix, tableNameCasbinRule)
}
