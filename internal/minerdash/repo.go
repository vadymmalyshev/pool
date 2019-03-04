package minerdash

import (
	"encoding/json"
	"fmt"

	"git.tor.ph/hiveon/pool/config"
	client "github.com/influxdata/influxdb1-client"
	. "github.com/influxdata/influxdb1-client/models"
	log "github.com/sirupsen/logrus"

	"strings"
)

type Repositorer interface {
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

type Repository struct {
	influx *client.Client
}

func NewRepository(client *client.Client) *Repository {
	return &Repository{client}
}

func (m *Repository) queryRaw(query string) interface{} {
	q := client.Query{
		Command:  query,
		Database: config.InfluxDB.Name,
	}

	res, err := m.influxClient.Query(q)
	if err != nil {
		log.Error(err)
	}

	if res == nil || len(res.Results < 0) {
		return nil
	}

	return res.Results[0]
}

func (m *Repository) queryFloatSingle(query string) (result float64) {
	res := m.query(q)

	if res != nil && res.Series != nil {
		result, _ = res.Series[0].Values[0][1].(json.Number).Float64()
	}

	return result
}

func (m *Repository) querySingle(query string) Row {
	res := m.query(q)

	if res != nil && res.Series != nil {
		return res.Series[0]
	}

	return Row{}
}

func (m *Repository) GetPoolLatestShare() Row {
	workerState := config.WorkerState

	sql := fmt.Sprintf(`
	SELECT mean(validShares) as validShares 
	FROM a_year.pool_hashrate 
	WHERE time>now()-2h AND time<now()-25m 
	GROUP BY time(%s) ORDER BY time DESC LIMIT 1`, workerState)

	return m.querySingle(sql)
}

func (m *Repository) GetPoolWorker() Row {
	sql := `
	SELECT max(count) as count 
	FROM a_year.worker_count 
	WHERE time>now()-1h`

	return m.querySingle(sql)
}

func (m *Repository) GetPoolMiner() Row {
	sql := `
	SELECT max(count) as count 
	FROM a_year.miner_count WHERE time>now()-1h`

	return m.querySingle(sql)
}

func (m *Repository) GetETHHashrate() IncomeCurrency {
	sql := fmt.Sprintf(`
	SELECT 
		mean(difficulty) as difficulty,
		mean(cny) as cny, 
		mean(usd) as usd
		mean(btc) as btc 
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

func (m *Repository) GetShares(walletID string, workerID string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom

	

	sql := fmt.Sprintf(`
	SELECT 
		sum(invalidShares) as invalidShares, 
		sum(validShares) as validShares, 
		sum(staleShares) as staleShares 
	FROM 
		two_hours.worker 
	WHERE wallet=%s  
		AND time>now()-%s `, walletID, time)

	// if workerID provided
	if len(strings.TrimSpace(workerId)) > 0 {
		sql += fmt.Sprintf(" AND rig='%s' ", workerID)
	}

	sql := fmt.Sprintf("GROUP BY time(%s) fill(0)", zoom)

	res := m.querySingle(sql)
	if res.Values != nil {
		return res
	}

	log.Error("Can't querySingle influx data")
	log.Info("Query: ", sql)

	return Row{}
}

func (m *Repository) GetLocalHashrateResult(walletID string, workerID string) Row {
	time := config.ZoomConfigTime
	zoom := config.ZoomConfigZoom

	sql := fmt.Sprintf(`
		SELECT sum(localHashrate) as localHashrate 
		FROM a_month.miner_worker 
		WHERE time>now()-%s`, time)

	// if workerID provided
	if len(strings.TrimSpace(workerId)) > 0 {
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

func (m *Repository) GetHashrate(walletID string, workerID string) RepoHashrate {
	sql := fmt.Sprintf(`
		SELECT mean(validShares) as validShares 
		FROM a_year.miner 
		WHERE wallet='%s' `, walletID)

	rigSQL = fmt.Sprintf(" AND rig='%s'", workerID)

	hashrateSQL := sql + " AND time>now()-1h AND time<now()-25m"
	meanHashRateSQL := sql + fmt.Sprintf(" AND token='' AND time>now()-%s AND time<now()-25m", config.PoolZoom))

	if workerId != "" {
		hashrateSQL += rigSQL
		meanHashRateSQL += rigSQL
	}

	return RepoHashrate{
		Hashrate: m.queryFloatSingle(hashrateSql), 
		Hashrate24H: m.queryFloatSingle(meanHashRateSql)
	}
}

func (m *Repository) GetCountHistory(walletID string) Row {
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

func (m *Repository) GetHashrate20m(walletID string, workerID string) client.Result {
	sql := fmt.Sprintf(`
		SELECT sum(validShares) as validShares,
		sum(invalidShares) as invalidShares,
		sum(staleShares) as staleShares,
		sum(originValidShares) as originValidShares 
		FROM two_hours.worker 
		WHERE wallet='%s' AND time > now()-20m and time <= now()`, walletID)

	if workerId != "" {
		sql += fmt.Sprintf(" and rig='%s'", workerID)
	}
	
	sql += " GROUP BY rig"

	return m.query(sql)
}

func (m *Repository) GetAvgHashrate1d(walletID string, workerID string) client.Result {
	sql := fmt.Sprintf(`
		SELECT mean(validShares) as validShares, 
		mean(invalidShares) as invalidShares, 
		mean(staleShares) as staleShares,
		mean(localHashrate) as localHashrate 
		FROM a_month.miner_worker 
		WHERE wallet='%s' and time>now()-%s`, walletID, config.PoolZoom)

	if len(strings.TrimSpace(workerId)) == 0 {
		sql += fmt.Sprintf(" and rig='%s'",workerID)
	}

	sql += " GROUP BY rig"

	return m.query(sql)
}

func (m *Repository) GetWorkers24hStatistic(walletID string, workerID string) client.Result {
	time := config.WorkerConfigTime
	zoom := config.WorkerConfigZoom
	sql := `
		SELECT sum(validShares) as validShares,
		sum(invalidShares) as invalidShares,
		sum(staleShares) as staleShares 
		FROM a_month.miner_worker 
		WHERE`

	if len(strings.TrimSpace(walletID)) > 0 {
		sql += fmt.Sprintf(" wallet='%s' and", walletId)
	}

	if len(strings.TrimSpace(workerID)) > 0 {
		sql += fmt.Sprintf(" rig='%s' and time>now()-%s group by time(%s), rig", workerID, time, zoom)
	} 
	
	sql += fmt.Sprintf(" time>now()-%s GROUP BY time(%s), rig", time, zoom)

	return m.query(sql)
}

func (m *Repository) GetWorkersWalletsMapping24hStatistic() client.Result {
	sql := fmt.Sprintf(`
		SELECT sum(validShares) as validShares,
		sum(invalidShares) as invalidShares,
		sum(staleShares) as staleShares 
		FROM a_month.miner_worker 
		WHERE time>now()-%s 
		GROUP BY rig, wallet`, config.WorkerConfigTime)
				
	return m.query(sql)
}
