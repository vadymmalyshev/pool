package middleware

import (
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/api/utils"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"net/http"
)

type Authorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	a := &Authorizer{enforcer: e}

	return func(c *gin.Context) {
		if config.UseCasbin {
			if !a.CheckPermission(c.Request) {
				a.RequirePermission(c)
			}
		}
		c.Next()
	}
}

func (a *Authorizer) CheckPermission(r *http.Request) bool {
	token := r.Context().Value("token").(swagger.OAuth2TokenIntrospection)
	user, _ := utils.GetUserByEmail(token.Sub)
	method := r.Method
	path := r.URL.Path
	return a.enforcer.Enforce(user, path, method)
}

func (a *Authorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}
