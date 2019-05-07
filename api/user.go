package api

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	"github.com/jinzhu/gorm"
	"strconv"
	"git.tor.ph/hiveon/pool/internal/users"
	"github.com/gin-gonic/gin"
)

const (
	paramFID      = "fid"
	paramCoin     = "coin"
	paramWalletID = "walletID"
)

type UserAPI struct {
	userService users.UserServicer
}

func NewUserAPI(admDB *gorm.DB) *UserAPI {
	return &UserAPI{userService: users.NewUserService(admDB)}
}

func (h *UserAPI) GetUserWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID, err := strconv.ParseUint(c.Param(paramWalletID), 10, 32)
		if apierrors.HandleError(err) {
			c.JSON(400, apierrors.NewApiErr(400, "Parse error"))
			return
		}
		wallet, err := h.userService.GetUserWallet(uint(walletID))
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, wallet)
	}
}

func (h *UserAPI) PostUserWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		fid, _ := strconv.ParseUint(c.PostForm(paramFID), 10, 32)
		walletID := c.PostForm(paramWalletID)
		coin := c.PostForm(paramCoin)

		_, err := h.userService.SaveUserWallet(uint(fid), walletID, coin)
		if apierrors.AbortWithApiError(c, err) {
			return
		}

		c.JSON(201, walletID)
	}
}
