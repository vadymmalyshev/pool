package admin

import (
	"git.tor.ph/hiveon/pool/internal/casbin"
	"regexp"

	"git.tor.ph/hiveon/pool/models"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
)

// SwitchDatabasesMiddleware used to switch databases on the fly
func SwitchDatabasesMiddleware(db, seq2, idp *gorm.DB) *admin.Middleware {
	return &admin.Middleware{
		Name: "switch_db",
		Handler: func(context *admin.Context, middleware *admin.Middleware) {
			// switch admin's database to db2 for products related requests
			/*
			if regexp.MustCompile(models.Wallet{}.AdminPath()).MatchString(context.Request.URL.Path) ||
				regexp.MustCompile(models.Coin{}.AdminPath()).MatchString(context.Request.URL.Path) {
				context.SetDB(db)
			}*/
			if regexp.MustCompile(models.Blacklist{}.AdminPath()).MatchString(context.Request.URL.Path) {
				context.SetDB(seq2)
			} else if regexp.MustCompile(casbin.CasbinRule{}.AdminPath()).MatchString(context.Request.URL.Path) {
				context.SetDB(idp)
			} else {
				context.SetDB(db)
			}
			middleware.Next(context)
		},
	}
}
