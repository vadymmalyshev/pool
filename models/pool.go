package models

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

// Success response
// swagger:response IncomeHistory
type IncomeHistory struct {
	Code int      `json:"code"`
	Data []Income `json:"data"`
}

type Income struct {
	Time   string `json:"time"`
	Income string `json:"income"`
}