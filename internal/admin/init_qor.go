package admin

import (
	"git.tor.ph/hiveon/idp/models/users"
	"git.tor.ph/hiveon/pool/internal/casbin"
	"git.tor.ph/hiveon/pool/models"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"net/http"
)

func InitAdmin(adminDB *gorm.DB, idpDB *gorm.DB, seq2 *gorm.DB, mux *http.ServeMux) {
	admin := admin.New(&admin.AdminConfig{DB: idpDB})
	admin.GetRouter().Use(SwitchDatabasesMiddleware(adminDB, seq2, idpDB))

	admin.AddResource(&models.Wallet{})
	admin.AddResource(&models.Coin{})
	admin.AddResource(&users.User{})
	admin.AddResource(&models.Blacklist{})
	admin.AddResource(&casbin.CasbinRule{})
	admin.AddResource(&models.WorkerFee{})
	admin.AddResource(&models.BillingWorkerStatistic{})
	admin.AddResource(&models.BillingWorkerMoney{})

	admin.MountTo("/admin", mux)
}
