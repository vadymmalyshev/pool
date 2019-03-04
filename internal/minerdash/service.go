package minerdash

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"time"

	"git.tor.ph/hiveon/pool/internal/income"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// Service provides method with calculations of miner's job history
type Service struct {
	db *gorm.DB
}

type service interface {
	GetFutureIncome() FutureIncome
	GetBillInfo(walletID string) BillInfo
	GetShares(walletID string, workerID string) Shares
	GetBill(walletID string) Bill
	GetMiner(walletID string, workerID string) MinerWorker
	GetHashrate(walletID string, workerID string) Hashrate
	GetCountHistory(walletID string) WorkerCount
	CalcWorkersStat(walletID string, workerID string) WorkersStatistic
	GetWalletWorkerMapping() WalletWorkerMappingStatistic
}

type minerService struct {
	income              income.IncomeRepositorer
	minerdashRepository IMinerdashRepository
	hiveosRepository    IHiveosRepository
	redisRepository     IRedisRepository
}

func NewMinerService() MinerService {
	return &minerService{blockRepository: NewBlockRepository(), minerdashRepository: NewRepository(),
		hiveosRepository: NewHiveosRepository(), redisRepository: NewRedisRepository()}
}

// for mockBlockRepo testing
func NewMinerServiceWithRepo(block IBlockRepository, minerdash IMinerdashRepository, hiveos IHiveosRepository, redis IRedisRepository) MinerService {
	return &minerService{blockRepository: block, minerdashRepository: minerdash, hiveosRepository: hiveos, redisRepository: redis}
}

func (m *minerService) GetFutureIncome() FutureIncome {
	incomeCurrency := m.minerdashRepository.GetETHHashrate()
	income := int(m.blockRepository.GetIncomeResult())
	income1d := int(m.blockRepository.GetIncome24h())

	FutureIncome{Code: 200}
	futureIncome.Data.CNY = incomeCurrency.CNY
	futureIncome.Data.BTC = incomeCurrency.BTC
	futureIncome.Data.USD = incomeCurrency.USD
	futureIncome.Data.Income = income
	futureIncome.Data.Income1d = income1d
	return futureIncome
}

func (m *minerService) GetBillInfo(walletID string) BillInfo {
	bill := m.hiveosRepository.GetBillInfo(walletID)
	return BillInfo{Balance: bill.Balance, FirstTime: FormatTimeToRFC3339(bill.FirstTime),
		FirstPaid: bill.FirstPaid, TotalPaid: bill.TotalPaid}
}

func (m *minerService) GetBill(walletID string) Bill {
	rows := m.hiveosRepository.GetBill(walletID)

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
		r.Time = FormatTimeToRFC3339(r.Time)
		bills = append(bills, r)
	}

	return Bill{Code: 200, Data: bills}
}

func (m *minerService) GetShares(walletID string, workerID string) Shares {
	rows := m.minerdashRepository.GetShares(walletID, workerID)

	sharesDetails := make([]SharesDetail, len(rows.Values))

	for i, row := range rows.Values {
		sharesDetails[i].Time = row[0].(string)
		sharesDetails[i].InvalidShares, _ = row[1].(json.Number).Float64()
		sharesDetails[i].ValidShares, _ = row[2].(json.Number).Float64()
		sharesDetails[i].StaleShares, _ = row[3].(json.Number).Float64()

	}

	removeExtraElements(&sharesDetails)
	calcMeanHashrate(&sharesDetails)

	rows = m.minerdashRepository.GetLocalHashrateResult(walletID, "")

	timeMap := make(map[string]float64, len(rows.Values))
	for _, row := range rows.Values {
		timeMap[row[0].(string)], _ = row[1].(json.Number).Float64()
	}

	for i := range sharesDetails {
		sharesDetails[i].LocalHashrate = timeMap[sharesDetails[i].Time]
	}

	return Shares{Code: 200, Data: sharesDetails}
}

func (m *minerService) GetMiner(walletID string, workerID string) MinerWorker {

	minerWorker := MinerWorker{}
	minerWorker.Balance = m.getBalance(walletID)
	minerWorker.Hashrate = m.GetHashrate(walletID, workerID)
	minerWorker.Workers = m.getLatestWorker(walletID)
	minerWorker.WorkerCounts = m.GetCountHistory(walletID)
	return minerWorker
}

func (m *minerService) getBalance(walletID string) Balance {

	balance := Balance{Code: 200}
	balance.Data.Balance = m.hiveosRepository.GetBalance(walletID)
	return balance
}

func (m *minerService) GetHashrate(walletID string, workerID string) Hashrate {

	hashrateRepo := m.minerdashRepository.GetHashrate(walletID, workerID)
	hashrate := Hashrate{Code: 200}
	hashrate.Data.Hashrate = calcHashrate(hashrateRepo.Hashrate)
	hashrate.Data.MeanHashrate24H = calcHashrate(hashrateRepo.Hashrate24H)
	return hashrate
}

//worker.list
func (m *minerService) getLatestWorker(walletID string) Workers {
	const msToSec = 1000000000
	workers2 := m.redisRepository.GetLatestWorker(walletID)
	workers := m.minerdashRepository.GetHashrate20m(walletID, "")
	workers1dHashrate := m.minerdashRepository.GetAvgHashrate1d(walletID, "")

	workersMap := make(map[string]Worker)
	for _, w := range workers.Series {
		worker := Worker{}
		worker.Rig = w.Tags["rig"]

		worker.Time = GetRowStringValue(w, 0, "time")
		worker.ValidShares = GetRowFloatValue(w, 0, "validShares")
		worker.InvalidShares = GetRowFloatValue(w, 0, "invalidShares")
		worker.StaleShares = GetRowFloatValue(w, 0, "staleShares")
		workersMap[worker.Rig] = worker
	}

	workers1dMap := make(map[string]Worker)
	for _, w := range workers1dHashrate.Series {
		worker := Worker{}
		worker.Rig = w.Tags["rig"]

		worker.Time = GetRowStringValue(w, 0, "time")
		worker.ValidShares = GetRowFloatValue(w, 0, "validShares")
		worker.InvalidShares = GetRowFloatValue(w, 0, "invalidShares")
		worker.MeanLocalHashrate1d = GetRowFloatValue(w, 0, "localHashrate")

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
		worker.Rig = k
		worker.Time = timeStamp.Format(time.RFC3339)
		worker.ValidShares = workersMap[k].ValidShares
		worker.StaleShares = workersMap[k].StaleShares
		worker.InvalidShares = workersMap[k].InvalidShares

		worker.Hashrate1d = calcHashrate(workers1dMap[k].ValidShares)
		worker.MeanLocalHashrate1d = workers1dMap[k].MeanLocalHashrate1d
		worker.Hashrate = calcHashrate(worker.ValidShares)
		workerResult = append(workerResult, worker)
	}
	return Workers{Code: 200, Data: workerResult}
}

func (m *minerService) GetCountHistory(walletID string) WorkerCount {
	row := m.minerdashRepository.GetCountHistory(walletID)
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
	return workerCount
}
func (m *minerService) GetWalletWorkerMapping() WalletWorkerMappingStatistic {
	var walletWorkerMappingList []WalletWorkerMapping
	mapping := m.minerdashRepository.GetWorkersWalletsMapping24hStatistic()
	for _, w := range mapping.Series {
		workerWallet := WalletWorkerMapping{}
		workerWallet.Wallet = w.Tags["wallet"]
		workerWallet.Worker = w.Tags["rig"]
		walletWorkerMappingList = append(walletWorkerMappingList, workerWallet)
	}
	return WalletWorkerMappingStatistic{200, walletWorkerMappingList}
}

func (m *minerService) CalcWorkersStat(walletID string, workerID string) WorkersStatistic {
	workers := m.minerdashRepository.GetWorkers24hStatistic(walletID, workerID)
	var workersList []Worker
	var workerStatisticList []WorkerStatistic

	for _, v := range workers.Series {
		for in, _ := range v.Values {
			worker := Worker{}
			worker.Rig = v.Tags["rig"]
			worker.Time = GetRowStringValue(v, in, "time")
			worker.ValidShares = GetRowFloatValue(v, in, "validShares")
			worker.InvalidShares = GetRowFloatValue(v, in, "invalidShares")
			worker.StaleShares = GetRowFloatValue(v, in, "staleShares")
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

	return WorkersStatistic{Code: 200, Data: workerStatisticList}
}

func createWorkerStatistic(v string, validSharesSum float64, invalidSharesSum float64, staleSharesSum float64,
	percentage float64, workerStatisticList *[]WorkerStatistic) {
	workerStat := WorkerStatistic{}
	workerStat.Rig = v
	workerStat.ValidShares = math.Round(validSharesSum)
	workerStat.InvalidShares = math.Round(invalidSharesSum)
	workerStat.StaleShares = math.Round(staleSharesSum)
	workerStat.ActivityPercentage = RoundFloat2(percentage)
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
	hashrate := GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCul") /
		GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCulDivider")

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

func calcHashrate(count float64) float64 {
	hashRateCul := GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCul") /
		GetConfig().GetFloat64("app.config.pool.hashrate.hashrateCulDivider")

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
