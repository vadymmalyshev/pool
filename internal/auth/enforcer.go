package auth

import "github.com/casbin/casbin"

const authModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.sub == p.sub || p.sub == "*") && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

// NewAccessEnforcer returns enforcer user by auth mw and synchronizer
func NewAccessEnforcer(a interface{}) *casbin.Enforcer {
	return casbin.NewEnforcer(casbin.NewModel(authModel), a)
}
