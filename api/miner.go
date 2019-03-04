package api

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/minerdash"
	"git.tor.ph/hiveon/pool/internal/platform/database/influx"
	"git.tor.ph/hiveon/pool/internal/platform/database/mysql"
	"github.com/gin-gonic/gin"
	red "github.com/go-redis/redis"
	"log"
	"strconv"
)

const (
	paramWorker = "workerId"
	paramWallet = "walletId"
)

type MinerAPI struct {
	minerService MinerServicer
}

func NewMinerAPI() *MinerAPI {

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

	return &MinerAPI{minerService: NewMinerService(Sequelize2DB, Sequelize3DB, client, red_client)}
}

func (h *MinerAPI) GetFutureIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.minerService.GetFutureIncome())
	}
}

func (h *MinerAPI) GetBillInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		c.JSON(200, h.minerService.GetBillInfo(walletID))
	}
}

func (h *MinerAPI) GetBill() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		c.JSON(200, h.minerService.GetBill(walletID))
	}
}

func (h *MinerAPI) GetShares() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		workerID := c.Param(paramWorker)
		c.JSON(200, h.minerService.GetShares(walletID, workerID))
	}
}

func (h *MinerAPI) GetHashrate() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		c.JSON(200, h.minerService.GetHashrate(walletID, ""))
	}
}

func (h *MinerAPI) GetCountHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		c.JSON(200, h.minerService.GetCountHistory(walletID))
	}
}

func (h *MinerAPI) GetWorkerStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		workerID := c.Param(paramWorker)
		c.JSON(200, h.minerService.CalcWorkersStat("", workerID))
	}
}

func (h *MinerAPI) GetAllWorkersStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.minerService.CalcWorkersStat("", ""))
	}
}

func (h *MinerAPI) GetWalletStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		c.JSON(200, h.minerService.CalcWorkersStat(walletID, ""))
	}
}

func (h *MinerAPI) GetWalletsWorkersMapping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.minerService.GetWalletWorkerMapping())
	}
}

func (h *MinerAPI) GetMiner() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		c.JSON(200, h.minerService.GetMiner(walletID, ""))
	}
}
