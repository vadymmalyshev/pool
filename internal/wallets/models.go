package wallets


type shareStatistic struct {
	reportedHashrate float64
	valid float64
	invalid float64
	reportedHashrate24h float64
	validShares24h float64
	invalidShares24h float64
	staleShares24h float64
	invalidSharesStake24h float64
	staleSharesStake24h float64
}

type expectedIncome struct {
	ETH1d float64
	ETH1dUSD float64
	ETH7d float64
	ETH7dUSD float64
}

type workerOnlineStatistic struct {
	online int
	offline int
	workerState  map[string]bool
}




