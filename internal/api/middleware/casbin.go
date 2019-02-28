package middleware

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/api/utils"
	"github.com/casbin/casbin"
	"net/http"
	. "github.com/ory/hydra/sdk/go/hydra/swagger"
	log "github.com/sirupsen/logrus"
)

func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			token := r.Context().Value("token").(OAuth2TokenIntrospection)
			if config.UseCasbin{

				method := r.Method
				path := r.URL.Path

				user, _ := GetUserByEmail(token.Sub)
				log.Info(user)

				e.LoadPolicy() //refresh policy
				if e.Enforce(fmt.Sprint(user.ID), path, method) {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, http.StatusText(403), 403)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
