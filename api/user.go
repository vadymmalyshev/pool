package api

import (
	. "git.tor.ph/hiveon/pool/internal/api/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserAPI struct {
	userService UserService
}

func NewUserAPI() *UserAPI {
	return &UserAPI{userService:NewUserService()}
}

func (h *UserAPI) GetUserWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID, _ := strconv.ParseUint(c.Param("walletId"), 10, 32)
		c.JSON(200, h.userService.GetUserWallet(uint(walletID)))
	}
}

func (h *UserAPI) PostUserWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		fid, _ := strconv.ParseUint(c.PostForm("fid"),10 ,32)
		walletID := c.PostForm("walletID")
		coin := c.PostForm("coin")
		h.userService.SaveUserWallet(uint(fid), walletID, coin)
		c.JSON(201, walletID)
	}
}