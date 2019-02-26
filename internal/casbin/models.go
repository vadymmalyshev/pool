package casbin

import (
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
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
