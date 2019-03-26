package minerdash

import (
	"encoding/json"
	"fmt"
	"git.tor.ph/hiveon/pool/api/apierrors"
	. "github.com/influxdata/influxdb1-client/models"

	"git.tor.ph/hiveon/pool/config"
	"github.com/influxdata/influxdb1-client"
	log "github.com/sirupsen/logrus"

	"strings"
)

type MinerdashRepositorer interface {
	GetPoolLatestShare() (Row, error)
	GetPoolWorker() (Row, error)
	GetPoolMiner() (Row, error)
	GetETHHashrate() (IncomeCurrency, error)
	GetShares(walletId string, workerId string) (Row, error)
	GetLocalHashrateResult(walletId string, workerId string) (Row, error)
	GetHashrate(walletId string, workerId string) (RepoHashrate, error)
	GetCountHistory(walletId string) (Row, error)
	GetHashrate20m(walletId string, workerId string) (client.Result, error)
	GetAvgHashrate1d(walletId string, workerId string) (client.Result, error)
	GetWorkers24hStatistic(walletId string, workerId string) (client.Result, error)
	GetWorkersWalletsMapping24hStatistic() (client.Result, error)
}

type MinerdashRepository struct {
	influx *client.Client
}

func NewMinerdashRepository(client *client.Client) *MinerdashRepository {
	return &MinerdashRepository{influx: client}
}

func (m *MinerdashRepository) queryRaw(query string) (interface{}, error) {
	q := client.Query{
		Command:  query,
		Database: config.InfluxDB.Name,
	}

	res, err := m.influx.Query(q)
	if apierrors.HandleError(err) {
		return nil, err
	}

	if len(res.Results) <= 0 {
		return nil, apierrors.NewApiErr(400, "No data")
	}

	return res.Results[0], nil
}

func (m *MinerdashRepository) queryFloatSingle(query string) (float64, error) {
	data, err := m.queryRaw(query)
	if err != nil {
		return 0, err
	}
	res := data.(client.Result)
	var result float64
	if res.Series == nil {
		return 0, apierrors.NewApiErr(400, "No data")
	}
	result, _ = res.Series[0].Values[0][1].(json.Number).Float64()

	return result, nil
}

func (m *MinerdashRepository) querySingle(query string) (Row, error) {
	data, err := m.queryRaw(query)
	if err != nil {
		return Row{}, err
	}
	res := data.(client.Result)
	if res.Series == nil {
		return Row{}, apierrors.NewApiErr(400, "No data")
	}

	return res.Series[0], nil
}

func (m *MinerdashRepository) GetPoolLatestShare() (Row, error) {
	workerState := config.WorkerState

	sql := fmt.Sprintf(`
	SELECT mean(validShares) as validShares 
	FROM a_year.pool_hashrate 
	WHERE time>now()-2h AND time<now()-25m 
	GROUP BY time(%s) ORDER BY time DESC LIMIT 1`, workerState)

	return m.querySingle(sql)
}

func (m *MinerdashRepository) GetPoolWorker() (Row, error) {
	sql := `
	SELECT max(count) as count 
	FROM a_year.worker_count 
	WHERE time>now()-1h`

	return m.querySingle(sql)
}

func (m *MinerdashRepository) GetPoolMiner() (Row, error) {
	sql := `
	SELECT max(count) as count 
	FROM a_year.miner_count WHERE time>now()-1h`

	return m.querySingle(sql)
}

func (m *MinerdashRepository) GetETHHashrate() (IncomeCurrency, error) {
	sql := fmt.Sprintf(`
	SELECT 
		mean(difficulty) as difficulty,
		mean(cny_float) as cny, 
		mean(usd)       as usd,
		mean(btc)       as btc 
	FROM a_year.eth_stats WHERE time>now()-%s`, config.PoolZoom)

	res, err := m.querySingle(sql)
	if err != nil {
		return IncomeCurrency{}, err
	}

	if res.Values == nil {
		return IncomeCurrency{}, apierrors.NewApiErr(400, "No data")
	}
	//TODO How to get data by column names?
	cny, _ := res.Values[0][2].(json.Number).Float64()
	usd, _ := res.Values[0][3].(json.Number).Float64()
	btc, _ := res.Values[0][4].(json.Number).Float64()
	return IncomeCurrency{CNY: cny, BTC: btc, USD: usd}, nil
}

func (m *MinerdashRepository) GetShares(walletID string, workerID string) (Row, error) {
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

	res, err := m.querySingle(sql)
	if err != nil {
		return Row{}, err
	}
	if res.Values == nil {
		log.Error("Can't querySingle influx data")
		log.Info("Query: ", sql)
		return Row{}, apierrors.NewApiErr(400, "No data")
	}

	return res, nil
}

func (m *MinerdashRepository) GetLocalHashrateResult(walletID string, workerID string) (Row, error) {
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

	res, err := m.querySingle(sql)
	if err != nil {
		return Row{}, err
	}
	if res.Values == nil {
		// seems like shit
		log.Error("Can't querySingle influx data")
		log.Info("Query: ", sql)
		return Row{}, apierrors.NewApiErr(400, "No data")
	}

	return res, nil
}

func (m *MinerdashRepository) GetHashrate(walletID string, workerID string) (RepoHashrate, error) {
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
	hashrate, err := m.queryFloatSingle(hashrateSQL)
	if err != nil {
		return RepoHashrate{}, err
	}
	hashrate24H, err := m.queryFloatSingle(meanHashRateSQL)
	if err != nil {
		return RepoHashrate{}, err
	}
	return RepoHashrate{
		Hashrate:    hashrate,
		Hashrate24H: hashrate24H,
	}, nil
}

func (m *MinerdashRepository) GetCountHistory(walletID string) (Row, error) {
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

	res, err := m.querySingle(sql)
	if err != nil {
		return Row{}, err
	}

	if res.Values == nil {
		log.Error("Can't querySingle influx data")
		log.Info("Query: ", sql)
		return Row{}, apierrors.NewApiErr(400, "No data")
	}
	return res, nil
}

func (m *MinerdashRepository) GetHashrate20m(walletID string, workerID string) (client.Result, error) {
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
	data, err := m.queryRaw(sql)
	if err != nil {
		return client.Result{}, err
	}

	return data.(client.Result), nil
}

func (m *MinerdashRepository) GetAvgHashrate1d(walletID string, workerID string) (client.Result, error) {
	sql := fmt.Sprintf(`
		SELECT mean(validShares) as validShares, 
		mean(invalidShares) as invalidShares, 
		mean(staleShares) as staleShares,
		mean(localHashrate) as localHashrate 
		FROM a_month.miner_worker 
		WHERE wallet='%s' and time>now()-%s`, walletID, config.PoolZoom)

	if len(strings.TrimSpace(workerID)) == 0 {
		sql += fmt.Sprintf(" and rig='%s'", workerID)
	}

	sql += " GROUP BY rig"
	data, err := m.queryRaw(sql)
	if err != nil {
		return client.Result{}, err
	}

	return data.(client.Result), nil
}

func (m *MinerdashRepository) GetWorkers24hStatistic(walletID string, workerID string) (client.Result, error) {
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
		sql += fmt.Sprintf(" rig='%s' and time>now()-%s and time<=now()-1h group by time(%s), rig", workerID, time, zoom)
	}

	sql += fmt.Sprintf(" time>now()-%s GROUP BY time(%s), rig", time, zoom)

	data, err := m.queryRaw(sql)
	if err != nil {
		return client.Result{}, err
	}

	return data.(client.Result), nil
}

func (m *MinerdashRepository) GetWorkersWalletsMapping24hStatistic() (client.Result, error) {
	sql := fmt.Sprintf(`
		SELECT sum(validShares) as validShares,
		sum(invalidShares) as invalidShares,
		sum(staleShares) as staleShares 
		FROM a_month.miner_worker 
		WHERE time>now()-%s 
		GROUP BY rig, wallet`, config.WorkerConfigTime)

	data, err := m.queryRaw(sql)
	if err != nil {
		return client.Result{}, err
	}

	return data.(client.Result), nil
}
