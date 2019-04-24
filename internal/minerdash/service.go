package minerdash

import (
	"encoding/json"
	"git.tor.ph/hiveon/pool/api/apierrors"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/accounting"
	"git.tor.ph/hiveon/pool/internal/api/utils"
	"git.tor.ph/hiveon/pool/internal/income"
	"git.tor.ph/hiveon/pool/internal/redis"
	influx "github.com/influxdata/influxdb1-client/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"sort"
	"strconv"
	"time"
)

// Service provides method with calculations of miner's job history
type Service struct {
	db *gorm.DB
}

type MinerServicer interface {
	GetFutureIncome() (FutureIncome, error)
	GetBillInfo(walletID string) (BillInfo, error)
	GetShares(walletID string, workerID string) (Shares, error)
	GetBill(walletID string) (Bill, error)
	GetMiner(walletID string, workerID string) (MinerWorker, error)
	GetHashrate(walletID string, workerID string) (Hashrate, error)
	GetCountHistory(walletID string) (WorkerCount, error)
	CalcWorkersStat(walletID string, workerID string) (WorkersStatistic, error)
	GetWalletWorkerMapping() (WalletWorkerMappingStatistic, error)
	CalcHashrate(count float64) float64
	GetIndex() (PoolData, error)
}

type minerService struct {
	incomeRepository     income.IncomeRepositorer
	minerdashRepository  MinerdashRepositorer
	accountingRepository accounting.AccointingRepositorer
	redisRepository      redis.RedisRepositorer
}

func NewMinerService() MinerServicer {
	return &minerService{incomeRepository: income.NewIncomeRepository(config.Seq3), minerdashRepository: NewMinerdashRepository(config.Influx),
		accountingRepository: accounting.NewAccountingRepository(config.Seq2), redisRepository: redis.NewRedisRepository(config.Red)}
}

// for mockBlockRepo testing
func NewMinerServiceWithRepo(incomeRepository income.IncomeRepositorer, minerdashRepository MinerdashRepositorer, accountingRepository accounting.AccointingRepositorer, redisRepository redis.RedisRepositorer) MinerServicer {
	return &minerService{incomeRepository: incomeRepository, minerdashRepository: minerdashRepository, accountingRepository: accountingRepository, redisRepository: redisRepository}
}

func (m *minerService) GetFutureIncome() (FutureIncome, error) {
	incomeCurrency, err := m.minerdashRepository.GetETHHashrate()
	if err != nil {
		return FutureIncome{}, err
	}
	income, err := m.incomeRepository.GetIncomeResult()
	if err != nil {
		return FutureIncome{}, err
	}
	income1d, err := m.incomeRepository.GetIncome24h()
	if err != nil {
		return FutureIncome{}, err
	}
	futureIncome := FutureIncome{Code: 200}
	futureIncome.Data.CNY = incomeCurrency.CNY
	futureIncome.Data.BTC = incomeCurrency.BTC
	futureIncome.Data.USD = incomeCurrency.USD
	futureIncome.Data.Income = int(income)
	futureIncome.Data.Income1d = int(income1d)
	return futureIncome, nil
}

func (m *minerService) GetBillInfo(walletId string) (BillInfo, error) {
	bill, err := m.accountingRepository.GetBillInfo(walletId)
	if err != nil {
		return BillInfo{}, err
	}
	return BillInfo{Balance: bill.Balance, FirstTime: utils.FormatTimeToRFC3339(bill.FirstTime),
		FirstPaid: bill.FirstPaid, TotalPaid: bill.TotalPaid}, nil
}

func (m *minerService) GetBill(walletId string) (Bill, error) {
	rows, err := m.accountingRepository.GetBill(walletId)
	if err != nil {
		return Bill{}, err
	}

	var bills []BillDetail

	for rows.Next() {
		var r BillDetail
		err := rows.Scan(&r.Id, &r.Paid, &r.Status, &r.Time, &r.TXHash)
		if err != nil {
			log.Error(err)
		}
		if r.Status == "9000" {
			r.Status = "SUCCESS"
		}
		r.Time = utils.FormatTimeToRFC3339(r.Time)
		bills = append(bills, r)
	}
	if bills == nil {
		return Bill{}, apierrors.NewApiErr(400, "Bad Request")
	}

	return Bill{Code: 200, Data: bills}, nil
}

func (m *minerService) GetShares(walletID string, workerID string) (Shares, error) {
	rows, err := m.minerdashRepository.GetShares(walletID, workerID)
	if err != nil {
		return Shares{}, err
	}

	sharesDetails := make([]SharesDetail, len(rows.Values))

	if len(rows.Values) == 0 {
		return Shares{}, apierrors.NewApiErr(400, "Bad request")
	}

	for i, row := range rows.Values {
		sharesDetails[i].Time = row[0].(string)
		sharesDetails[i].InvalidShares, _ = row[1].(json.Number).Float64()
		sharesDetails[i].ValidShares, _ = row[2].(json.Number).Float64()
		sharesDetails[i].StaleShares, _ = row[3].(json.Number).Float64()

	}

	removeExtraElements(&sharesDetails)
	calcMeanHashrate(&sharesDetails)

	rows, err = m.minerdashRepository.GetLocalHashrateResult(walletID, "")
	if err != nil {
		return Shares{}, err
	}

	timeMap := make(map[string]float64, len(rows.Values))
	for _, row := range rows.Values {
		timeMap[row[0].(string)], _ = row[1].(json.Number).Float64()
	}

	for i := range sharesDetails {
		sharesDetails[i].LocalHashrate = timeMap[sharesDetails[i].Time]
	}

	return Shares{Code: 200, Data: sharesDetails}, nil
}

func (m *minerService) GetMiner(walletID string, workerID string) (MinerWorker, error) {

	minerWorker := MinerWorker{}
	balance, err := m.getBalance(walletID)
	if err != nil {
		return MinerWorker{}, err
	}
	minerWorker.Balance = balance
	minerWorker.Hashrate, err = m.GetHashrate(walletID, workerID)
	if err != nil {
		return MinerWorker{}, err
	}
	minerWorker.Workers, err = m.getLatestWorker(walletID)
	if err != nil {
		return MinerWorker{}, err
	}
	minerWorker.WorkerCounts, err = m.GetCountHistory(walletID)
	if err != nil {
		return MinerWorker{}, err
	}
	return minerWorker, nil
}

func (m *minerService) getBalance(walletID string) (Balance, error) {

	balance := Balance{Code: 200}
	value, err := m.accountingRepository.GetBalance(walletID)
	if err != nil {
		return Balance{}, err
	}
	balance.Data.Balance = value
	return balance, nil
}

func (m *minerService) GetHashrate(walletID string, workerID string) (Hashrate, error) {

	hashrateRepo, err := m.minerdashRepository.GetHashrate(walletID, workerID)
	if err != nil {
		return Hashrate{}, err
	}
	hashrate := Hashrate{Code: 200}
	hashrate.Data.Hashrate = m.CalcHashrate(hashrateRepo.Hashrate)
	hashrate.Data.MeanHashrate24H = m.CalcHashrate(hashrateRepo.Hashrate24H)
	return hashrate, nil
}

func (m *minerService) GetIndex() (PoolData, error) {
	hashrateCul, _ := strconv.ParseFloat(config.HashrateCul, 64)
	hashrateCulDivider, _ := strconv.ParseFloat(config.HashrateCulDivider, 64)
	hashrateConfig := hashrateCul / hashrateCulDivider

	hashRate, err := m.minerdashRepository.GetPoolLatestShare()
	if err != nil {
		return PoolData{}, err
	}
	miner, err := m.minerdashRepository.GetPoolMiner()
	if err != nil {
		return PoolData{}, err
	}
	worker, err := m.minerdashRepository.GetPoolWorker()
	if err != nil {
		return PoolData{}, err
	}

	poolData := PoolData{Code: 200}

	if !reflect.DeepEqual(hashRate, influx.Row{}) {
		poolData.Data.Hashrate.Time = hashRate.Values[0][0].(string)
		validShares := hashRate.Values[0][1]
		if validShares != nil {
			poolData.Data.Hashrate.ValidShares, _ = validShares.(json.Number).Float64()
		}

		if poolData.Data.Hashrate.ValidShares > 0 {
			val := math.Round(poolData.Data.Hashrate.ValidShares * hashrateConfig)
			if math.IsNaN(val) {
				val = 0
			}
			poolData.Data.Hashrate.Hashrate = val
		}
	}

	if !reflect.DeepEqual(miner, influx.Row{}) {
		poolData.Data.Miner.Time = miner.Values[0][0].(string)
		minerCount, _ := miner.Values[0][1].(json.Number).Float64()
		val := math.Round(minerCount/1000*10) / 10
		if math.IsNaN(val) {
			val = 0
		}
		poolData.Data.Miner.Count = val
	}

	if !reflect.DeepEqual(worker, influx.Row{}) {
		poolData.Data.Worker.Time = worker.Values[0][0].(string)
		workerCount, _ := worker.Values[0][1].(json.Number).Float64()
		val := math.Round(workerCount/1000*10) / 10
		if math.IsNaN(val) {
			val = 0
		}
		poolData.Data.Worker.Count = val
	}

	return poolData, nil
}

//worker.list
func (m *minerService) getLatestWorker(walletID string) (Workers, error) {
	const msToSec = 1000000000
	workers2, err := m.redisRepository.GetLatestWorker(walletID)
	if err != nil {
		return Workers{}, err
	}
	workers, err := m.minerdashRepository.GetHashrate20m(walletID, "")
	if err != nil {
		return Workers{}, err
	}
	workers1dHashrate, err := m.minerdashRepository.GetAvgHashrate1d(walletID, "")
	if err != nil {
		return Workers{}, err
	}

	workersMap := make(map[string]Worker)
	for _, w := range workers.Series {
		worker := Worker{}
		worker.Rig = utils.FormatWorkerName(w.Tags["rig"])
		worker.Time = utils.GetRowStringValue(w, 0, "time")
		worker.ValidShares = utils.GetRowFloatValue(w, 0, "validShares")
		worker.InvalidShares = utils.GetRowFloatValue(w, 0, "invalidShares")
		worker.StaleShares = utils.GetRowFloatValue(w, 0, "staleShares")
		workersMap[worker.Rig] = worker
	}

	workers1dMap := make(map[string]Worker)
	for _, w := range workers1dHashrate.Series {
		worker := Worker{}
		worker.Rig = utils.FormatWorkerName(w.Tags["rig"])
		worker.Time = utils.GetRowStringValue(w, 0, "time")
		worker.ValidShares = utils.GetRowFloatValue(w, 0, "validShares")
		worker.InvalidShares = utils.GetRowFloatValue(w, 0, "invalidShares")
		worker.MeanLocalHashrate1d = utils.GetRowFloatValue(w, 0, "localHashrate")

		workers1dMap[worker.Rig] = worker
	}

	var workerResult []Worker

	for k, v := range workers2 {
		ts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Error(err)
			continue
		}
		timeStamp := time.Unix(ts/msToSec, 0)
		if timeStamp.Before(time.Now().Add(-time.Duration(24) * time.Hour)) {
			continue
		}

		worker := Worker{}
		worker.Rig = utils.FormatWorkerName(k)
		worker.Time = timeStamp.Format(time.RFC3339)
		worker.ValidShares = workersMap[k].ValidShares
		worker.StaleShares = workersMap[k].StaleShares
		worker.InvalidShares = workersMap[k].InvalidShares

		worker.Hashrate1d = m.CalcHashrate(workers1dMap[k].ValidShares)
		worker.MeanLocalHashrate1d = workers1dMap[k].MeanLocalHashrate1d
		worker.Hashrate = m.CalcHashrate(worker.ValidShares)
		workerResult = append(workerResult, worker)
	}
	return Workers{Code: 200, Data: workerResult}, nil
}

func (m *minerService) GetCountHistory(walletID string) (WorkerCount, error) {
	row, err := m.minerdashRepository.GetCountHistory(walletID)
	if err != nil {
		return WorkerCount{}, err
	}
	var timeCount []TimeCount

	if row.Values != nil {
		for _, el := range row.Values {
			tc := TimeCount{}
			tc.Time = el[0].(string)
			tc.Count, _ = el[1].(json.Number).Float64()
			timeCount = append(timeCount, tc)
		}
	}

	if len(timeCount) > 1 {
		if timeCount[0].Count == 0 {
			timeCount = timeCount[1:]
		}
		if timeCount[len(timeCount)-1].Count == 0 {
			timeCount = timeCount[:len(timeCount)-1]
		}
	}
	workerCount := WorkerCount{Code: 200}
	workerCount.Data = timeCount
	return workerCount, nil
}
func (m *minerService) GetWalletWorkerMapping() (WalletWorkerMappingStatistic, error) {
	var walletWorkerMappingList []WalletWorkerMapping
	mapping, err := m.minerdashRepository.GetWorkersWalletsMapping24hStatistic()
	if err != nil {
		return WalletWorkerMappingStatistic{}, err
	}
	for _, w := range mapping.Series {
		workerWallet := WalletWorkerMapping{}
		workerWallet.Wallet = w.Tags["wallet"]
		workerWallet.Worker = w.Tags["rig"]
		walletWorkerMappingList = append(walletWorkerMappingList, workerWallet)
	}
	return WalletWorkerMappingStatistic{200, walletWorkerMappingList}, nil
}

func (m *minerService) CalcWorkersStat(walletID string, workerID string) (WorkersStatistic, error) {
	workers, err := m.minerdashRepository.GetWorkers24hStatistic(walletID, workerID)
	if err != nil {
		return WorkersStatistic{}, err
	}
	var workersList []Worker
	var workerStatisticList []WorkerStatistic

	for _, v := range workers.Series {
		for in, _ := range v.Values {
			worker := Worker{}
			worker.Rig = utils.FormatWorkerName(v.Tags["rig"])
			worker.Time = utils.GetRowStringValue(v, in, "time")
			worker.ValidShares = utils.GetRowFloatValue(v, in, "validShares")
			worker.InvalidShares = utils.GetRowFloatValue(v, in, "invalidShares")
			worker.StaleShares = utils.GetRowFloatValue(v, in, "staleShares")
			workersList = append(workersList, worker)
		}
	}
	sort.Sort(RigSorter(workersList))
	currentWorker := workersList[0].Rig
	validSharesSum := 0.00
	invalidSharesSum := 0.00
	staleSharesSum := 0.00
	percentage := 0.00

	for _, v := range workersList {
		if currentWorker != v.Rig { // new time series
			createWorkerStatistic(currentWorker, validSharesSum, invalidSharesSum, staleSharesSum, percentage, &workerStatisticList)

			validSharesSum = 0.00
			invalidSharesSum = 0.00
			staleSharesSum = 0.00
			percentage = 0.00
			currentWorker = v.Rig
		}
		validSharesSum += v.ValidShares
		invalidSharesSum += v.InvalidShares
		staleSharesSum += v.StaleShares

		if v.ValidShares > 0 || v.InvalidShares > 0 || v.StaleShares > 0 {
			percentage += 0.3472 // 5min %
		}
	}
	createWorkerStatistic(currentWorker, validSharesSum, invalidSharesSum, staleSharesSum, percentage, &workerStatisticList) //last worker

	return WorkersStatistic{Code: 200, Data: workerStatisticList}, nil
}

func createWorkerStatistic(v string, validSharesSum float64, invalidSharesSum float64, staleSharesSum float64,
	percentage float64, workerStatisticList *[]WorkerStatistic) {
	workerStat := WorkerStatistic{}
	workerStat.Rig = v
	workerStat.ValidShares = math.Round(validSharesSum)
	workerStat.InvalidShares = math.Round(invalidSharesSum)
	workerStat.StaleShares = math.Round(staleSharesSum)
	workerStat.ActivityPercentage = utils.RoundFloat2(percentage)
	*workerStatisticList = append(*workerStatisticList, workerStat)
}

//this func to remove last element if time > time.now() - 20 mins
func removeExtraElements(sharesDetails *[]SharesDetail) {
	if len(*sharesDetails) > 1 {
		*sharesDetails = (*sharesDetails)[1:]
	}

	if len(*sharesDetails) > 1 {
		timeStamp, _ := time.Parse(time.RFC3339, (*sharesDetails)[len(*sharesDetails)-1].Time)
		ago20 := time.Now().Add(-time.Millisecond * 1200000)

		if timeStamp.After(ago20) {
			*sharesDetails = (*sharesDetails)[:len(*sharesDetails)-1]
		}
	}
}

func calcMeanHashrate(sharesDetails *[]SharesDetail) {
	var (
		count    float64
		numCount float64
	)
	beginIsZero := true
	hashrate := utils.GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCul") /
		utils.GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCulDivider")

	for i := range *sharesDetails {
		if !beginIsZero || (*sharesDetails)[i].ValidShares > 0 || (*sharesDetails)[i].InvalidShares > 0 {
			numCount++
			beginIsZero = false
		}

		if (*sharesDetails)[i].ValidShares > 0 {
			(*sharesDetails)[i].Hashrate = math.Round(((*sharesDetails)[i].ValidShares) * hashrate)
			count += (*sharesDetails)[i].Hashrate
		}

		if beginIsZero {
			(*sharesDetails)[i].MeanHashrate = 0
		} else {
			(*sharesDetails)[i].MeanHashrate = math.Round(count/numCount*100) / 100
		}
	}
}

func (m *minerService) CalcHashrate(count float64) float64 {
	hashRateCul := utils.GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCul") /
		utils.GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCulDivider")

	result := math.Round(hashRateCul * count)
	if math.IsNaN(result) {
		return 0
	}
	return result
}

// sort

type RigSorter []Worker

func (a RigSorter) Len() int           { return len(a) }
func (a RigSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RigSorter) Less(i, j int) bool { return a[i].Rig < a[j].Rig }
