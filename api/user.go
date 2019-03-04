package api

import (
	"strconv"

	. "git.tor.ph/hiveon/pool/internal/api/service"
	"github.com/gin-gonic/gin"
)

const (
	paramFID      = "fid"
	paramCoin     = "coin"
	paramWalletID = "walletID"
)

type UserAPI struct {
	userService UserService
}

func NewUserAPI() *UserAPI {
	return &UserAPI{userService: NewUserService()}
}

func (h *UserAPI) GetUserWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID, _ := strconv.ParseUint(c.Param(paramWalletID), 10, 32)
		c.JSON(200, h.userService.GetUserWallet(uint(walletID)))
	}
}

func (h *UserAPI) PostUserWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		fid, _ := strconv.ParseUint(c.PostForm(paramFID), 10, 32)
		walletID := c.PostForm(paramWalletID)
		coin := c.PostForm(paramCoin)

		h.userService.SaveUserWallet(uint(fid), walletID, coin)

		c.JSON(201, walletID)
	}
}
