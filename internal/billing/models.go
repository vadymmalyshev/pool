package billing

// Success response
// swagger:response WalletEarning
type WalletEarning struct {
	Code int     `json:"code"`
	Data Earning `json:"data"`
}

type Earning struct {
	Address        string  `json:"address"`
	Date           string  `json:"date"`
	Hashrate       float64 `json:"hashrate"`
	USD            float64 `json:"usd"`
	CNY            float64 `json:"cny"`
	BTC            float64 `json:"btc"`
	Commission_USD float64 `json:"commission_usd"`
}
