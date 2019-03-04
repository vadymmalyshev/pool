package minerdash

type IncomeCurrency struct {
	CNY float64
	BTC float64
	USD float64
}

type RepoHashrate struct {
	Hashrate    float64
	Hashrate24H float64
}

// Success response
// swagger:response FutureIncome
type FutureIncome struct {
	Code int `json:"code"`
	Data struct {
		CNY      float64 `json:"cny"`
		USD      float64 `json:"usd"`
		BTC      float64 `json:"btc"`
		Income1d int     `json:"income1d"`
		Income   int     `json:"income"`
	} `json:"data"`
}

// Success response
// swagger:response BillInfo
type BillInfo struct {
	Balance   string  `json:"balance"`
	FirstPaid float64 `json:"firstPaid"`
	FirstTime string  `json:"firstTime"`
	TotalPaid float64 `json:"totalPaid"`
}

// Success response
// swagger:response Shares
type Shares struct {
	Code int            `json:"code"`
	Data []SharesDetail `json:"data"`
}

type SharesDetail struct {
	Hashrate      float64 `json:"hashrate"`
	InvalidShares float64 `json:"invalidShares"`
	LocalHashrate float64 `json:"localHashrate"`
	MeanHashrate  float64 `json:"meanHashrate"`
	StaleShares   float64 `json:"staleShares"`
	Time          string  `json:"time"`
	ValidShares   float64 `json:"validShares"`
}

// Success response
// swagger:response Bill
type Bill struct {
	Code int          `json:"code"`
	Data []BillDetail `json:"data"`
}

type BillDetail struct {
	Id     int    `json:"-"`
	Paid   string `json:"paid"`
	Status string `json:"status"`
	TXHash string `json:"tx_hash"`
	Time   string `json:"time"`
}

//Success response
// swagger:response MinerWorker
type MinerWorker struct {
	Balance      Balance     `json:"balance"`
	WorkerCounts WorkerCount `json:"workerCounts"`
	Hashrate     Hashrate    `json:"hashrate"`
	Workers      Workers     `json:"workers"`
}

// Success response
// swagger:response Balance
type Balance struct {
	Code int `json:"code"`
	Data struct {
		Balance float64 `json:"balance"`
	} `json:"data"`
}

// Success response
// swagger:response WorkerCount
type WorkerCount struct {
	Code int         `json:"code"`
	Data []TimeCount `json:"data"`
}

type TimeCount struct {
	Time  string  `json:"time"`
	Count float64 `json:"count"`
}

// Success response
// swagger:response Hashrate
type Hashrate struct {
	Code int `json:"code"`
	Data struct {
		Hashrate        float64 `json:"hashrate"`
		MeanHashrate24H float64 `json:"meanHashrate24H"`
	} `json:"data"`
}

type Workers struct {
	Code int      `json:"code"`
	Data []Worker `json:"data"`
}

type Worker struct {
	Rig                 string  `json:"rig"`
	Time                string  `json:"time"`
	Hashrate1d          float64 `json:"hashrate1d"`
	MeanLocalHashrate1d float64 `json:"meanLocalHashrate1d"`
	InvalidShares       float64 `json:"invalidShares"`
	StaleShares         float64 `json:"staleShares"`
	ValidShares         float64 `json:"validShares"`
	Hashrate            float64 `json:"hashrate"`
}

type WorkersStatistic struct {
	Code int               `json:"code"`
	Data []WorkerStatistic `json:"data"`
}

type WalletWorkerMappingStatistic struct {
	Code int                   `json:"code"`
	Data []WalletWorkerMapping `json:"data"`
}

type WorkerStatistic struct {
	Rig                string  `json:"rig"`
	InvalidShares      float64 `json:"invalidShares"`
	StaleShares        float64 `json:"staleShares"`
	ValidShares        float64 `json:"validShares"`
	ActivityPercentage float64 `json:"activityPercentage"`
}

type WalletWorkerMapping struct {
	Worker string `json:"worker"`
	Wallet string `json:"wallet"`
}
