package repository

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb1-client"
	"git.tor.ph/hiveon/pool/config"
	influx "git.tor.ph/hiveon/pool/internal/platform/database/influx"
	. "github.com/influxdata/influxdb1-client/models"
	log "github.com/sirupsen/logrus"

	"strings"
)

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

type IncomeCurrency struct {
	CNY float64
	BTC float64
	USD float64
}

type RepoHashrate struct {
	Hashrate    float64
	Hashrate24H float64
}

func NewMinerdashRepository() *MinerdashRepository {
	return &MinerdashRepository{influxClient: GetMinerdashClient()}
}

func GetMinerdashClient() *client.Client {
	client, _ := influx.Connect(config.InfluxDB)
	return client
}

func (m *MinerdashRepository) queryFloatSingle(query string) (result float64) {
	q := client.Query{
		Command:  query,
		Database: config.InfluxDB.Name,
	}

	res, err := m.influxClient.Query(q)
	if err != nil {
		log.Error(err)
	}

	if res.Results[0].Series != nil {
		result, _ = res.Results[0].Series[0].Values[0][1].(json.Number).Float64()
	}

	return result
}

func (m *MinerdashRepository) querySingle(query string) Row {
	q := client.Query{
		Command:  query,
		Database: config.InfluxDB.Name,
	}

	res, err := m.influxClient.Query(q)
	if err != nil {
		log.Error(err)
	}

	if res.Results != nil && res.Results[0].Series != nil {
		return res.Results[0].Series[0]
	}

	return Row{}
}

func (m *MinerdashRepository) query(query string) client.Result {
	q := client.Query{
		Command:  query,
		Database: config.InfluxDB.Name,
	}

	res, err := m.influxClient.Query(q)
	if err != nil {
		log.Error(err)
	}
	return res.Results[0]
}

func (m *MinerdashRepository) GetPoolLatestShare() Row {
	workerState := config.WorkerState
	latestShareSQL := fmt.Sprintf("select mean(validShares) as validShares from a_year.pool_hashrate "+
		"where time>now()-2h and time<now()-25m group by time(%s) order by time desc limit 1", workerState)
	return m.querySingle(latestShareSQL)
}

func (m *MinerdashRepository) GetPoolWorker() Row {
	minerSQL := "select max(count) as count from a_year.worker_count where time>now()-1h"
	return m.querySingle(minerSQL)
}

func (m *MinerdashRepository) GetPoolMiner() Row {
	workerSQL := "select max(count) as count from a_year.miner_count where time>now()-1h"
	return m.querySingle(workerSQL)
}


func (m *MinerdashRepository) GetETHHashrate() IncomeCurrency {
	sql := fmt.Sprintf("select mean(difficulty) as difficulty,mean(cny_float) as cny, mean(usd) as usd, "+
		"mean(btc) as btc from a_year.eth_stats where time>now()-%s", config.PoolZoom)

	res := m.querySingle(sql)

	//TODO How to get data by column names?
	if res.Values != nil {
		cny, _ := res.Values[0][2].(json.Number).Float64()
		usd, _ := res.Values[0][3].(json.Number).Float64()
		btc, _ := res.Values[0][4].(json.Number).Float64()
		return IncomeCurrency{CNY: cny, BTC: btc, USD: usd}
	}
	return IncomeCurrency{}
}

func (m *MinerdashRepository) GetShares(walletId string, workerId string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom
	sql := "select sum(invalidShares) as invalidShares, sum(validShares) as validShares, sum(staleShares) "+
		"as staleShares from two_hours.worker where "
	if len(strings.TrimSpace(workerId)) == 0 {
		sql = fmt.Sprintf(sql + "wallet='%s' and time>now()-%s group by time(%s) fill(0)", walletId, time, zoom)
	} else {
		sql = fmt.Sprintf(sql + "wallet='%s' and rig='%s' and time>now()-%s group by time(%s) fill(0)", walletId, workerId, time, zoom)
	}
	res := m.querySingle(sql)
	if res.Values != nil {
		return res
	}
	{
		log.Error("Can't querySingle influx data")
		log.Info("Query: ", sql)
		return Row{}
	}

}

func (m *MinerdashRepository) GetLocalHashrateResult(walletId string, workerId string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom
	sql := "select sum(localHashrate) as localHashrate from a_month.miner_worker where "
	if workerId == "" {
		sql = fmt.Sprintf(sql +  "wallet='%s' "+
			"and time>now()-%s group by time(%s) fill(0)", walletId, time, zoom)
	} else {
		sql = fmt.Sprintf(sql +  "wallet='%s' "+
			"and rig='%s' and time>now()-%s group by time(%s) fill(0)", walletId, workerId, time, zoom)
	}

	res := m.querySingle(sql)
	if res.Values != nil {
		return res
	}
	{
		log.Error("Can't querySingle influx data")
		log.Info("Query: ", sql)
		return Row{}
	}

}

func (m *MinerdashRepository) GetHashrate(walletId string, workerId string) RepoHashrate {
	sql := "select mean(validShares) as validShares from a_year.miner where "
	var hashrateSql, meanHashRateSql string

	if  workerId != "" {
		hashrateSql = fmt.Sprintf(sql + "wallet='%s' and rig='%s' and time>now()-1h and time<now()-25m", walletId, workerId)
		meanHashRateSql = fmt.Sprintf(sql + "wallet='%s' and token='' and rig='%s' and time>now()-%s and time<now()-25m", walletId, workerId, config.PoolZoom)
	} else {
		hashrateSql = fmt.Sprintf(sql + "wallet='%s' and time>now()-1h and time<now()-25m", walletId)
		meanHashRateSql = fmt.Sprintf(sql + "wallet='%s' and token='' and time>now()-%s and time<now()-25m", walletId, config.PoolZoom)
	}

	return RepoHashrate{Hashrate: m.queryFloatSingle(hashrateSql), Hashrate24H: m.queryFloatSingle(meanHashRateSql)}
}

func (m *MinerdashRepository) GetCountHistory(walletId string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom
	var sql string
	if len(time) > 0 {
		sql = fmt.Sprintf("select count(a) as count from (select count(validShares) as a from a_month.miner_worker "+
			"where time>now()-%s and wallet='%s' group by time(%s),rig) where time>now()-%s "+
			"group by time(%s) fill(0)", time, walletId, zoom, time, zoom)
	} else {
		sql = fmt.Sprintf("select max(count) as count from a_month.miner_worker_count where"+
			" time>now()-%s and wallet='%s' group by time(%s)", time, walletId, zoom)
	}

	res := m.querySingle(sql)
	if res.Values != nil {
		return res
	}
	{
		log.Error("Can't querySingle influx data")
		log.Info("Query: ", sql)
		return Row{}
	}

}

func (m *MinerdashRepository) GetHashrate20m(walletId string, workerId string) client.Result {
	sql := "select sum(validShares) as validShares,sum(invalidShares) as invalidShares,sum(staleShares)"+
		" as staleShares,sum(originValidShares) as originValidShares from two_hours.worker where "

	if workerId == "" {
		sql = fmt.Sprintf(sql + "wallet='%s' and time > now()-20m and time <= now() group by rig", walletId)
	} else {
		sql = fmt.Sprintf(sql + "wallet='%s' and rig='%s' and time > now()-20m and time <= now() group by rig", walletId, workerId)
	}

	return m.query(sql)
}

func (m *MinerdashRepository) GetAvgHashrate1d(walletId string, workerId string) client.Result {
	sql := "select mean(validShares) as validShares, mean(invalidShares) as invalidShares, mean(staleShares) as staleShares,"+
		" mean(localHashrate) as localHashrate from a_month.miner_worker where "

	if len(strings.TrimSpace(workerId)) == 0 {
		sql = fmt.Sprintf(sql + "wallet='%s' and time>now()-%s group by rig", walletId, config.PoolZoom)
	} else {
		sql = fmt.Sprintf(sql + "wallet='%s' and rig='%s' and time>now()-%s group by rig", walletId, workerId, config.PoolZoom)
	}

	return m.query(sql)
}

func (m *MinerdashRepository) GetWorkers24hStatistic(walletId string, workerId string) client.Result {
	time := config.WorkerConfigTime
	zoom := config.WorkerConfigZoom
	sql := "select sum(validShares) as validShares,sum(invalidShares) as invalidShares,sum(staleShares)"+
		" as staleShares from a_month.miner_worker where "
	if len(strings.TrimSpace(walletId)) > 0 {
		sql = fmt.Sprintf(sql+"wallet='%s' and ", walletId)
	}
	if len(strings.TrimSpace(workerId)) == 0 {
		sql = fmt.Sprintf(sql + "time>now()-%s group by time(%s), rig", time, zoom)
	} else {
		sql = fmt.Sprintf(sql + "rig='%s' and time>now()-%s group by time(%s), rig", workerId, time, zoom)
	}

	return m.query(sql)
}

func (m *MinerdashRepository) GetWorkersWalletsMapping24hStatistic() client.Result {
	sql := "select sum(validShares) as validShares,sum(invalidShares) as invalidShares,sum(staleShares)"+
		" as staleShares from a_month.miner_worker where "
	sql = fmt.Sprintf(sql + "time>now()-%s group by rig, wallet", config.WorkerConfigTime)

	return m.query(sql)
}


