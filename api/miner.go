package api

import (
	"git.tor.ph/hiveon/pool/api/apierrors"
	. "git.tor.ph/hiveon/pool/internal/minerdash"
	"github.com/gin-gonic/gin"
)

const (
	paramWorker = "workerID"
	paramWallet = "walletID"
	paramDate   = "date"
)

type MinerAPI struct {
	minerService MinerServicer
}

func NewMinerAPI() *MinerAPI {
	return &MinerAPI{minerService: NewMinerService()}
}

// Handle GET /api/pool/futureIncome
func (h *MinerAPI) GetFutureIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		futureInc, err := h.minerService.GetFutureIncome()
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, futureInc)
	}
}

// Handle GET /api/eth/:walletID/billInfo
func (h *MinerAPI) GetBillInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		billInfo, err := h.minerService.GetBillInfo(walletID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, billInfo)
	}
}

// Handle GET /api/eth/:walletID/bill
func (h *MinerAPI) GetBill() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		bill, err := h.minerService.GetBill(walletID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, bill)
	}
}

func (h *MinerAPI) GetShares() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		workerID := c.Param(paramWorker)
		shares, err := h.minerService.GetShares(walletID, workerID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, shares)
	}
}

func (h *MinerAPI) GetHashrate() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		hashRate, err := h.minerService.GetHashrate(walletID, "")
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, hashRate)
	}
}

func (h *MinerAPI) GetCountHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		history, err := h.minerService.GetCountHistory(walletID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, history)
	}
}

func (h *MinerAPI) GetWorkerStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		workerID := c.Param(paramWorker)
		wStat, err := h.minerService.CalcWorkersStat("", workerID)
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, wStat)
	}
}

func (h *MinerAPI) GetAllWorkersStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		wStat, err := h.minerService.CalcWorkersStat("", "")
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, wStat)
	}
}

func (h *MinerAPI) GetWalletStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		wStat, err := h.minerService.CalcWorkersStat(walletID, "")
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, wStat)
	}
}

func (h *MinerAPI) GetWalletsWorkersMapping() gin.HandlerFunc {
	return func(c *gin.Context) {
		wWorMap, err := h.minerService.GetWalletWorkerMapping()
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, wWorMap)
	}
}

// Handle GET /api/eth/:walletID/workers/list
func (h *MinerAPI) GetMiner() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param(paramWallet)
		miner, err := h.minerService.GetMiner(walletID, "")
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, miner)
	}
}

func (h *MinerAPI) HandleGetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		index, err := h.minerService.GetIndex()
		if apierrors.AbortWithApiError(c, err) {
			return
		}
		c.JSON(200, index)
	}
}
