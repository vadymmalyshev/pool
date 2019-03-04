package middleware

import (
	"git.tor.ph/hiveon/pool/config"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/ory/hydra/sdk/go/hydra/swagger"
	. "git.tor.ph/hiveon/pool/internal/api/utils"
)
type Authorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	a := &Authorizer{enforcer: e}

	return func(c *gin.Context) {
		if config.UseCasbin{
			if (!a.CheckPermission(c.Request)) {
				a.RequirePermission(c)
			}
		}
		c.Next()
	}
}

func (a *Authorizer) CheckPermission(r *http.Request) bool {
	token := r.Context().Value("token").(OAuth2TokenIntrospection)
	user, _ := GetUserByEmail(token.Sub)
	method := r.Method
	path := r.URL.Path
	return a.enforcer.Enforce(user, path, method)
}

func (a *Authorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}


