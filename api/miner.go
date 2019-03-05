package api

import (
	. "git.tor.ph/hiveon/pool/internal/minerdash"
	"github.com/gin-gonic/gin"
)

const (
	paramWorker = "workerId"
	paramWallet = "walletId"
)

type MinerAPI struct {
	minerService MinerServicer
}

func NewMinerAPI() *MinerAPI {
	return &MinerAPI{minerService: NewMinerService()}
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

func (h *MinerAPI) HandleGetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, h.minerService.GetIndex())
	}
}
