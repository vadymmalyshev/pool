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

type WalletInfo struct {
	Code    int            `json:"code"`
	Total   WalletTotal    `json:"total"`
	Shares  []SharesDetail `json:"shares"`
	Workers []WalletWorker `json:"workers"`
	History []TimeCount    `json:"history"`
	Payouts []BillDetail   `json:"payouts"`
}

type WalletTotal struct {
	Hashrate            float64 `json:"hashrate"`
	MeanHashrate        float64 `json:"meanHashrate"`
	ReportedHashrate    float64 `json:"reportedHashrate"`
	ReportedHashrate24h float64 `json:"reportedHashrate24h"`
	Valid               float64 `json:"valid"`
	Invalid             float64 `json:"invalid"`
	Balance             float64 `json:"balance"`
	Valid24h            float64 `json:"valid24h"`
	Stale24h            float64 `json:"stale24h"`
	Invalid24h          float64 `json:"invalid24h"`
	Stale24hStake       float64 `json:"stale24hStake"`
	Invalid24hStake     float64 `json:"invalid24hStake"`
	Expected24hUSD      float64 `json:"expected24hUSD"`
	Expected24h         float64 `json:"expected24h"`
	Expected7d          float64 `json:"expected7d"`
	Expected7dUSD       float64 `json:"expected7dUSD"`
	Online              int     `json:"online"`
	Offline             int     `json:"offline"`
}

type WalletWorker struct {
	Worker          string  `json:"worker"`
	Time            string  `json:"time"`
	Hashrate        float64 `json:"hashrate"`
	Hashrate24h     float64 `json:"hashrate24h"`
	MeanHashrate24h float64 `json:"meanHashrate24h"`
	Invalid         float64 `json:"invalid"`
	Stale           float64 `json:"stale"`
	Valid           float64 `json:"valid"`
	Online          bool    `json:"online"`
}

type WorkerInfo struct {
	Code   int            `json:"code"`
	Total  WorkerTotal    `json:"total"`
	Shares []SharesDetail `json:"shares"`
}

type WorkerTotal struct {
	Hashrate            float64 `json:"hashrate"`
	MeanHashrate        float64 `json:"meanHashrate"`
	ReportedHashrate    float64 `json:"reportedHashrate"`
	ReportedHashrate24h float64 `json:"reportedHashrate24h"`
	Valid               float64 `json:"valid"`
	Invalid             float64 `json:"invalid"`
	Valid24h            float64 `json:"valid24h"`
	Stale24h            float64 `json:"stale24h"`
	Invalid24h          float64 `json:"invalid24h"`
	Stale24hStake       float64 `json:"stale24hStake"`
	Invalid24hStake     float64 `json:"invalid24hStake"`
}

type PoolData struct {
	Code int `json:"code"`
	Data struct {
		Hashrate struct {
			Time        string  `json:"time"`
			ValidShares float64 `json:"validShares"`
			Hashrate    float64 `json:"hashrate"`
		} `json:"hashrate"`
		Miner struct {
			Time  string  `json:"time"`
			Count float64 `json:"count"`
		} `json:"miner"`
		Worker struct {
			Time  string  `json:"time"`
			Count float64 `json:"count"`
		} `json:"worker"`
	} `json:"data"`
}
