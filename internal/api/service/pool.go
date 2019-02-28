package service

import (
	"encoding/json"
	"git.tor.ph/hiveon/pool/config"
	. "github.com/influxdata/influxdb1-client/models"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"time"
	"reflect"
	. "git.tor.ph/hiveon/pool/internal/api/response"
	. "git.tor.ph/hiveon/pool/internal/api/repository"
)

type PoolService interface {
	GetIndex() PoolData
	GetIncomeHistory() IncomeHistory
}

type poolerService struct {
	minerdashRepository IMinerdashRepository
	blockRepository     IBlockRepository
}

func NewPoolService() PoolService {
	return &poolerService{NewMinerdashRepository(), NewBlockRepository()}
}

//for mockBlockRepo testing
func NewPoolServiceWithRepo(minerdashRepository IMinerdashRepository) PoolService {
	return &poolerService{minerdashRepository: minerdashRepository}
}

func (p *poolerService) GetIncomeHistory() IncomeHistory {
	rows := p.blockRepository.GetIncomeHistory()
	var incomeSlice []Income

	for rows.Next() {
		var income Income
		var t int64
		err := rows.Scan(&t, &income.Income)
		if err != nil {
			log.Error(err)
		}
		income.Time = time.Unix(t, 0).Format(time.RFC3339)
		incomeSlice = append(incomeSlice, income)
	}

	incomeHistory := IncomeHistory{Code: 200}
	incomeHistory.Data = incomeSlice
	return incomeHistory
}

func (p *poolerService) GetIndex() PoolData {
	hashrateCul, _ := strconv.ParseFloat(config.HashrateCul, 64)
	hashrateCulDivider, _ := strconv.ParseFloat(config.HashrateCulDivider, 64)
	hashrateConfig := hashrateCul / hashrateCulDivider

	hashRate := p.minerdashRepository.GetPoolLatestShare()
	miner := p.minerdashRepository.GetPoolMiner()
	worker := p.minerdashRepository.GetPoolWorker()

	poolData := PoolData{Code: 200}

	if !reflect.DeepEqual(hashRate, Row{}) {
		poolData.Data.Hashrate.Time = hashRate.Values[0][0].(string)
		validShares := hashRate.Values[0][1]
		if validShares != nil {
			poolData.Data.Hashrate.ValidShares, _ = validShares.(json.Number).Float64()
		}

		if poolData.Data.Hashrate.ValidShares > 0 {
			val := math.Round(poolData.Data.Hashrate.ValidShares * hashrateConfig)
			if math.IsNaN(val) {val = 0}
			poolData.Data.Hashrate.Hashrate = val
		}
	}

	if !reflect.DeepEqual(miner, Row{}) {
		poolData.Data.Miner.Time = miner.Values[0][0].(string)
		minerCount, _ := miner.Values[0][1].(json.Number).Float64()
		val := math.Round(minerCount/1000*10) / 10
		if math.IsNaN(val) {val = 0}
		poolData.Data.Miner.Count = val
	}

	if !reflect.DeepEqual(worker, Row{}) {
		poolData.Data.Worker.Time = worker.Values[0][0].(string)
		workerCount, _ := worker.Values[0][1].(json.Number).Float64()
		val := math.Round(workerCount/1000*10) / 10
		if math.IsNaN(val) {val = 0}
		poolData.Data.Worker.Count = val
	}

	return poolData
}
