package minerdash

import (
	"encoding/json"
	"fmt"
	. "github.com/influxdata/influxdb1-client/models"

	"git.tor.ph/hiveon/pool/config"
	client "github.com/influxdata/influxdb1-client"
	log "github.com/sirupsen/logrus"

	"strings"
)

type MinerdashRepositorer interface {
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
	influx *client.Client
}

func NewMinerdashRepository(client *client.Client) *MinerdashRepository {
	return &MinerdashRepository{influx : client}
}

func (m *MinerdashRepository) queryRaw(query string) interface{} {
	q := client.Query{
		Command:  query,
		Database: config.InfluxDB.Name,
	}

	res, err := m.influx.Query(q)
	if err != nil {
		log.Error(err)
	}

	if len(res.Results) <= 0 {
		return nil
	}

	return res.Results[0]
}

func (m *MinerdashRepository) queryFloatSingle(query string) (result float64) {
	res := m.queryRaw(query).(client.Result)

	if res.Series != nil {
		result, _ = res.Series[0].Values[0][1].(json.Number).Float64()
	}

	return result
}

func (m *MinerdashRepository) querySingle(query string) Row {
	res := m.queryRaw(query).(client.Result)

	if res.Series != nil {
		return res.Series[0]
	}

	return Row{}
}

func (m *MinerdashRepository) GetPoolLatestShare() Row {
	workerState := config.WorkerState

	sql := fmt.Sprintf(`
	SELECT mean(validShares) as validShares 
	FROM a_year.pool_hashrate 
	WHERE time>now()-2h AND time<now()-25m 
	GROUP BY time(%s) ORDER BY time DESC LIMIT 1`, workerState)

	return m.querySingle(sql)
}

func (m *MinerdashRepository) GetPoolWorker() Row {
	sql := `
	SELECT max(count) as count 
	FROM a_year.worker_count 
	WHERE time>now()-1h`

	return m.querySingle(sql)
}

func (m *MinerdashRepository) GetPoolMiner() Row {
	sql := `
	SELECT max(count) as count 
	FROM a_year.miner_count WHERE time>now()-1h`

	return m.querySingle(sql)
}

func (m *MinerdashRepository) GetETHHashrate() IncomeCurrency {
	sql := fmt.Sprintf(`
	SELECT 
		mean(difficulty) as difficulty,
		mean(cny_float) as cny, 
		mean(usd)       as usd,
		mean(btc)       as btc 
	FROM a_year.eth_stats WHERE time>now()-%s`, config.PoolZoom)

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

func (m *MinerdashRepository) GetShares(walletID string, workerID string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom

	

	sql := fmt.Sprintf(`
	SELECT 
		sum(invalidShares) as invalidShares, 
		sum(validShares) as validShares, 
		sum(staleShares) as staleShares 
	FROM 
		two_hours.worker 
	WHERE wallet='%s'
		AND time>now()-%s `, walletID, time)

	// if workerID provided
	if len(strings.TrimSpace(workerID)) > 0 {
		sql += fmt.Sprintf(" AND rig='%s' ", workerID)
	}

	sql += fmt.Sprintf(" GROUP BY time(%s) fill(0)", zoom)

	res := m.querySingle(sql)
	if res.Values != nil {
		return res
	}

	log.Error("Can't querySingle influx data")
	log.Info("Query: ", sql)

	return Row{}
}

func (m *MinerdashRepository) GetLocalHashrateResult(walletID string, workerID string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom

	sql := fmt.Sprintf(`
		SELECT sum(localHashrate) as localHashrate 
		FROM a_month.miner_worker 
		WHERE time>now()-%s`, time)

	// if workerID provided
	if len(strings.TrimSpace(workerID)) > 0 {
		sql += fmt.Sprintf(" AND rig='%s'", workerID)
	}

	sql += fmt.Sprintf(" GROUP BY time(%s) fill(0) ", zoom)
	
	res := m.querySingle(sql)
	if res.Values != nil {
		return res
	}

	// seems like shit
	log.Error("Can't querySingle influx data")
	log.Info("Query: ", sql)
	return Row{}
}

func (m *MinerdashRepository) GetHashrate(walletID string, workerID string) RepoHashrate {
	sql := fmt.Sprintf(`
		SELECT mean(validShares) as validShares 
		FROM a_year.miner 
		WHERE wallet='%s' `, walletID)

	rigSQL := fmt.Sprintf(" AND rig='%s'", workerID)

	hashrateSQL := sql + " AND time>now()-1h AND time<now()-25m"
	meanHashRateSQL := sql + fmt.Sprintf(" AND token='' AND time>now()-%s AND time<now()-25m", config.PoolZoom)

	if workerID != "" {
		hashrateSQL += rigSQL
		meanHashRateSQL += rigSQL
	}

	return RepoHashrate{
		Hashrate: m.queryFloatSingle(hashrateSQL),
		Hashrate24H: m.queryFloatSingle(meanHashRateSQL),
	}
}

func (m *MinerdashRepository) GetCountHistory(walletID string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom
	var sql string

	// what these scripts do??
	if len(time) > 0 {
		sql = fmt.Sprintf(`
		SELECT count(a) as count 
		FROM 
			(select count(validShares) as a 
			FROM a_month.miner_worker
			WHERE time>now()-%s AND wallet='%s' 
			GROUP BY time(%s),rig) 
		WHERE time>now()-%s 
		GROUP BY time(%s) fill(0)`, time, walletID, zoom, time, zoom)
	} else {
		sql = fmt.Sprintf(
			`SELECT max(count) as count 
			FROM a_month.miner_worker_count 
			WHERE time>now()-%s AND wallet='%s' 
			GROUP BY time(%s)`, time, walletID, zoom)
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

func (m *MinerdashRepository) GetHashrate20m(walletID string, workerID string) client.Result {
	sql := fmt.Sprintf(`
		SELECT sum(validShares) as validShares,
		sum(invalidShares) as invalidShares,
		sum(staleShares) as staleShares,
		sum(originValidShares) as originValidShares 
		FROM two_hours.worker 
		WHERE wallet='%s' AND time > now()-20m and time <= now()`, walletID)

	if workerID != "" {
		sql += fmt.Sprintf(" and rig='%s'", workerID)
	}
	
	sql += " GROUP BY rig"

	return m.queryRaw(sql).(client.Result)
}

func (m *MinerdashRepository) GetAvgHashrate1d(walletID string, workerID string) client.Result {
	sql := fmt.Sprintf(`
		SELECT mean(validShares) as validShares, 
		mean(invalidShares) as invalidShares, 
		mean(staleShares) as staleShares,
		mean(localHashrate) as localHashrate 
		FROM a_month.miner_worker 
		WHERE wallet='%s' and time>now()-%s`, walletID, config.PoolZoom)

	if len(strings.TrimSpace(workerID)) == 0 {
		sql += fmt.Sprintf(" and rig='%s'",workerID)
	}

	sql += " GROUP BY rig"

	return m.queryRaw(sql).(client.Result)
}

func (m *MinerdashRepository) GetWorkers24hStatistic(walletID string, workerID string) client.Result {
	time := config.WorkerConfigTime
	zoom := config.WorkerConfigZoom
	sql := `
		SELECT sum(validShares) as validShares,
		sum(invalidShares) as invalidShares,
		sum(staleShares) as staleShares 
		FROM a_month.miner_worker 
		WHERE`

	if len(strings.TrimSpace(walletID)) > 0 {
		sql += fmt.Sprintf(" wallet='%s' and", walletID)
	}

	if len(strings.TrimSpace(workerID)) > 0 {
		sql += fmt.Sprintf(" rig='%s' and time>now()-%s group by time(%s), rig", workerID, time, zoom)
	} 
	
	sql += fmt.Sprintf(" time>now()-%s GROUP BY time(%s), rig", time, zoom)

	return m.queryRaw(sql).(client.Result)
}

func (m *MinerdashRepository) GetWorkersWalletsMapping24hStatistic() client.Result {
	sql := fmt.Sprintf(`
		SELECT sum(validShares) as validShares,
		sum(invalidShares) as invalidShares,
		sum(staleShares) as staleShares 
		FROM a_month.miner_worker 
		WHERE time>now()-%s 
		GROUP BY rig, wallet`, config.WorkerConfigTime)
				
	return m.queryRaw(sql).(client.Result)
}
