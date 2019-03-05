package wallets

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/income"
	"git.tor.ph/hiveon/pool/internal/minerdash"
	. "git.tor.ph/hiveon/pool/internal/redis"
	"encoding/json"
	red "github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"math"
	"strconv"
	"time"
)

// Config stores connections to databases
type Config struct {
	Redis *red.Conn
	DB    *gorm.DB
}

type WalletServicer interface {
	GetWalletInfo(walletId string) minerdash.WalletInfo
	GetWalletWorkerInfo(walletId string, workerId string) minerdash.WorkerInfo
}

type walletService struct {
	minerService        minerdash.MinerServicer
	redisRepository     RedisRepositorer
	incomeRepository    IncomeRepositorer
	minerdashRepository minerdash.MinerdashRepositorer
}

func NewWalletService() WalletServicer {
	return &walletService{minerService: minerdash.NewMinerService(), redisRepository: NewRedisRepository(config.Red),
		incomeRepository: NewIncomeRepository(config.Seq3), minerdashRepository: minerdash.NewMinerdashRepository(config.Influx)}
}

func (w *walletService) GetWalletInfo(walletId string) minerdash.WalletInfo {
	miner := w.minerService.GetMiner(walletId,"")
	hashRates := miner.Hashrate
	balance := miner.Balance.Data.Balance
	workers := miner.Workers.Data
	history := miner.WorkerCounts.Data

	shares := w.minerService.GetShares(walletId, "").Data
	shareStat := w.getShareStatistic(shares)

	payouts := w.minerService.GetBill(walletId).Data
	workerStat := w.getWorkersStatistic(walletId)
	newWorkers := w.makeNewWorkers(workers, workerStat)

	futureIncomeData := w.minerService.GetFutureIncome().Data
	income1d := float64(futureIncomeData.Income1d)
	usd := futureIncomeData.USD
	income7d := w.incomeRepository.GetIncome7d()
	expectedIncome := w.getExpectedIncome(workers, income1d, income7d, usd)

	walletTotal := minerdash.WalletTotal{Hashrate: hashRates.Data.Hashrate, MeanHashrate: hashRates.Data.MeanHashrate24H, ReportedHashrate: shareStat.reportedHashrate,
		ReportedHashrate24h: shareStat.reportedHashrate24h, Valid: shareStat.valid, Invalid: shareStat.invalid, Balance: balance, Valid24h: shareStat.validShares24h, Stale24h: shareStat.staleShares24h,
		Invalid24h: shareStat.invalidShares24h, Stale24hStake: shareStat.staleSharesStake24h, Invalid24hStake: shareStat.invalidSharesStake24h, Expected24hUSD: expectedIncome.ETH1dUSD, Expected24h: expectedIncome.ETH1d,
		Expected7d: expectedIncome.ETH7d, Expected7dUSD: expectedIncome.ETH7dUSD, Online: workerStat.online, Offline: workerStat.offline}

	walletInfo := minerdash.WalletInfo{Code: 200, Total: walletTotal, Shares: shares, Workers: newWorkers, History: history, Payouts: payouts}

	return walletInfo
}

func (w *walletService) GetWalletWorkerInfo(walletId string, workerId string) minerdash.WorkerInfo {
	shares := w.minerService.GetShares(walletId, workerId).Data
	shareStat := w.getShareStatistic(shares)
	workers1dHashrate := w.minerdashRepository.GetAvgHashrate1d(walletId, workerId)

	var worker1d minerdash.Worker

	for _, w := range workers1dHashrate.Series {
		worker1d = minerdash.Worker{}
		worker1d.Rig = w.Tags["rig"]
		worker1d.Time = w.Values[0][0].(string)
		worker1d.ValidShares, _ = w.Values[0][1].(json.Number).Float64()
		worker1d.InvalidShares, _ = w.Values[0][2].(json.Number).Float64()
		worker1d.MeanLocalHashrate1d, _ = w.Values[0][4].(json.Number).Float64()
	}

	workerTotal := minerdash.WorkerTotal{Hashrate: w.minerService.CalcHashrate(worker1d.ValidShares), MeanHashrate: worker1d.MeanLocalHashrate1d, ReportedHashrate: shareStat.reportedHashrate,
		ReportedHashrate24h: shareStat.reportedHashrate24h, Valid: shareStat.valid, Invalid: shareStat.invalid, Valid24h: shareStat.validShares24h, Stale24h: shareStat.staleShares24h,
		Invalid24h: shareStat.invalidShares24h, Stale24hStake: shareStat.staleSharesStake24h, Invalid24hStake: shareStat.invalidSharesStake24h}

	workerInfo := minerdash.WorkerInfo{Code: 200, Total: workerTotal, Shares: shares}
	return workerInfo
}

// if there is activity in the last 20 minutes - online; calculate online and offline workers
func (w *walletService) getWorkersStatistic(walletId string) workerOnlineStatistic {
	const msToSec = 1000000000
	online := 0
	offline := 0
	workersState := make(map[string]bool)

	workers := w.redisRepository.GetLatestWorker(walletId)
	for k, v := range workers {
		ts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Error(err)
		}
		timeStamp := time.Unix(ts/msToSec, 0)
		if timeStamp.After(time.Now().Add(-time.Duration(20) * time.Minute)) {
			workersState[k] = true
			online ++
		} else {
			workersState[k] = false
			offline ++
		}
	}
	return workerOnlineStatistic{online, offline, workersState}
}

// get shares statistic
func (w *walletService) getShareStatistic(shares []minerdash.SharesDetail) shareStatistic {
	var reportedHashrate24h, reportedHashrate float64
	var validShares24h, invalidShares24h, staleShares24h float64 // all shares sum  24h
	var Stale24hStake, Invalid24hStake float64 // percent of all shares 24h
	var valid, invalid float64

	size := float64(len(shares))
	for i := range shares {
		reportedHashrate24h = reportedHashrate24h + shares[i].LocalHashrate
		validShares24h += shares[i].ValidShares
		invalidShares24h += shares[i].InvalidShares
		staleShares24h +=  shares[i].StaleShares
	}

	if size > 0 {
		reportedHashrate = shares[len(shares)-1].LocalHashrate
		valid = shares[len(shares)-1].ValidShares
		invalid = shares[len(shares)-1].InvalidShares
		reportedHashrate24h = reportedHashrate24h / size
		// stake
		totalShares := validShares24h + invalidShares24h + staleShares24h
		Stale24hStake = staleShares24h/totalShares
		Invalid24hStake = invalidShares24h/totalShares
	}
	return shareStatistic{reportedHashrate, valid, invalid, reportedHashrate24h,
		math.Round(validShares24h), math.Round(invalidShares24h), math.Round(staleShares24h),
		math.Round(Invalid24hStake), math.Round(Stale24hStake)}
}

// calculate expected income
func (w *walletService) getExpectedIncome(workers []minerdash.Worker, futureIncome1d float64, futureIncome7d float64, usd float64) expectedIncome {
	var res float64

	for k := range workers {
		res = res + workers[k].Hashrate1d
	}

	resETH1d := res * futureIncome1d
	resETH1dUSD := resETH1d * usd
	resETH7d := res * futureIncome7d
	resETH7dUSD := resETH7d * usd

	return expectedIncome{resETH1d, resETH1dUSD, resETH7d, resETH7dUSD}
}

// convert to new_worker format and check online status
func (w *walletService) makeNewWorkers(workers []minerdash.Worker, stat workerOnlineStatistic) []minerdash.WalletWorker {
	var res []minerdash.WalletWorker

	for _, v := range workers {
		workerNew := minerdash.WalletWorker{}
		workerNew.Worker = v.Rig
		workerNew.Hashrate = v.Hashrate
		workerNew.Valid = v.ValidShares
		workerNew.Hashrate24h = v.Hashrate1d
		workerNew.Invalid = v.InvalidShares
		workerNew.MeanHashrate24h = v.MeanLocalHashrate1d
		workerNew.Stale = v.StaleShares
		workerNew.Time = v.Time
		workerNew.Online = stat.workerState[v.Rig]

		res = append(res, workerNew)
	}
	return res
}
/*
// GetWorkersPulse return workers pulse map.
func (w WalletService) GetWorkersPulse() (map[string]WorkerPulse, error) {
	return GetWorkersPulse(*w.s.Redis, w.Address)
}
*/