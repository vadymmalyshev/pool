package api

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/wallets"
	"git.tor.ph/hiveon/pool/internal/platform/database/influx"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	red "github.com/go-redis/redis"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type WalletAPI struct {
	walletService WalletServicer
}

func NewWalletAPI() *WalletAPI {
	// TODO: in config
	Sequelize2DB, err := mysql.Connect(config.Sequelize2DB)
	if err != nil {
		log.Panic("failed to init mysql Sequelize2DB db :", err.Error())
	}

	Sequelize3DB, err := mysql.Connect(config.Sequelize3DB)
	if err != nil {
		log.Panic("failed to init mysql Sequelize3DB db :", err.Error())
	}

	client, err := influx.Connect(config.InfluxDB)
	if err != nil {
		log.Panic("failed to init influx:", err.Error())
	}

	DBName, _ := strconv.Atoi(config.Redis.Name)
	red_client := red.NewClient(&red.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: config.Redis.Pass,
		DB:       DBName,
	})
	_, err = red_client.Ping().Result()
	if err != nil {
		log.Panic("failed to init redis:", err.Error())
	}

	return &WalletAPI{walletService:NewWalletService(Sequelize2DB, Sequelize3DB, client, red_client)}
}

func (h *WalletAPI) HandleGetWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletId")
		c.JSON(200, h.walletService.GetWalletInfo(walletID))
	}
}

func (h *WalletAPI) HandleGetWalletWorker() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletId")
		workerID := c.Param("workerId")
		c.JSON(200, h.walletService.GetWalletWorkerInfo(walletID, workerID))
	}
}
