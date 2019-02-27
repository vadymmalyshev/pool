package minerdash

import client "github.com/influxdata/influxdb1-client"

type IMinerdashRepository interface {
	GetPoolLatestShare() Row
	GetPoolWorker() Row
	GetPoolMiner() Row
	GetETHHashrate() IncomeCurrency
	GetShares(walletId string, workerId string) Row
	GetLocalHashrateResult(walletId string, workerId string) Row
	GetHashrate(walletId string, workerId string) RepoHashrate
	GetCountHistory(walletId string) Row
	GetHashrate20m(walletId string, workerId string) client.Result
	GetAvgHashrate1d(walletId string, workerId string) client.Result
	GetWorkers24hStatistic(walletId string, workerId string) client.Result
	GetWorkersWalletsMapping24hStatistic() client.Result
}

type MinerdashRepository struct {
	influxClient *client.Client
}
