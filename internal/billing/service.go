package billing

import (
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/billing/utils"
	"git.tor.ph/hiveon/pool/models"
	"github.com/mileusna/crontab"
	"github.com/sirupsen/logrus"
)

type BillingCalculator struct {
	BillingRepo *BillingRepository
}

func NewBillingCalculator() *BillingCalculator {
	return &BillingCalculator{BillingRepo: NewBillingRepository()}
}

func (b BillingCalculator) StartCalculation(er chan error) {
	cTab := crontab.New()
	if err := cTab.AddJob("0 1 * * *", b.loadWorkerWalletStatistic); err != nil {
		logrus.Error("Failed to start billing module: ", err)
		er <- err
		return
	}
}

func (b BillingCalculator) loadWorkerWalletStatistic() {
	rates := b.fetchCurrencyRates()
	WalletWorkerMapping := b.consumeMapping()
	if err := b.generateStatistic(WalletWorkerMapping, rates); err != nil {
		logrus.Error(err)
	}
}

func (b BillingCalculator) fetchCurrencyRates() map[string]interface{} {
	logrus.Info("Consuming the currency rates : ", time.Now())
	ethAPI := "http://127.0.0.1:8090/api/pool/futureIncome" //testing
	resp, err := http.Get(ethAPI)
	if err != nil {
		logrus.Error(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	return utils.ParseJSON(string(body), false)
}

func (b BillingCalculator) consumeMapping() map[string]string {
	logrus.Info("Consuming wallet/worker mapping started: ", time.Now())

	MappingAPI := config.MappingApi
	resp, err := http.Get(MappingAPI)
	if err != nil {
		logrus.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}

	WalletWorkerMapping := make(map[string]string)
	res := utils.ParseJSON(string(body), false)

	for k, v := range res {
		if k == "data" {
			value := v.([]interface{})
			for _, v1 := range value {
				value1 := v1.(map[string]interface{})
				worker := value1["worker"].(string)
				wallet := value1["wallet"].(string)
				if worker != "" && wallet != "" {
					WalletWorkerMapping[worker] = wallet
				}
			}
		}
	}
	logrus.Info("Consuming wallet/worker mapping finished: ", time.Now())
	return WalletWorkerMapping
}

func (b BillingCalculator) generateStatistic(WalletWorkerMapping map[string]string,
	rates map[string]interface{}) error {
	logrus.Info("Consuming wallet/worker statistic started: ", time.Now())

	WorkersAPI := config.WorkersAPI
	resp, err := http.Get(WorkersAPI)
	if err != nil {
		logrus.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}
	res := utils.ParseJSON(string(body), false)

	// hashrate and currency rates
	hashrateCul, _ := strconv.ParseFloat(config.HashrateCul, 64)
	hashrateCulDivider, _ := strconv.ParseFloat(config.HashrateCulDivider, 64)
	hashrateConfig := hashrateCul / hashrateCulDivider
	USD := rates["usd"].(float64)
	BTC := rates["btc"].(float64)
	CNY := rates["cny"].(float64)

	currencyMap := make(map[string]float64)
	currencyMap["usd"] = USD
	currencyMap["btc"] = BTC
	currencyMap["cny"] = CNY

	for k, v := range res {
		if k == "data" {
			value := v.([]interface{})
			for _, v1 := range value {
				value1 := v1.(map[string]interface{})

				worker := value1["rig"].(string)
				validShares := value1["validShares"].(float64)
				invalidShares := value1["invalidShares"].(float64)
				staleShares := value1["staleShares"].(float64)
				percentage := value1["activityPercentage"].(float64)

				if worker != "" {
					wallet := WalletWorkerMapping[worker]
					stat := models.BillingWorkerStatistic{ValidShares: validShares, InvalidShares: invalidShares,
						StaleShares: staleShares, ActivityPercentage: percentage}
					work, err := b.BillingRepo.SaveWorkerStatistic(stat, wallet, worker)

					if err != nil {
						logrus.Error(err)
					} else {
						return b.calculateAndSaveCommission(stat, hashrateConfig, currencyMap, work)
					}
				}
			}
		}
	}
	logrus.Info("Consuming wallet/worker statistic finished: ", time.Now())

	return nil
}

func (b BillingCalculator) calculateAndSaveCommission(stat models.BillingWorkerStatistic, hashrateConfig float64,
	rates map[string]float64, worker *models.Worker) error {
	hashrate := stat.ValidShares * hashrateConfig
	hashrate_ := hashrate / 100000000
	USD := roundFloat(hashrate_ * rates["usd"])
	BTC := roundFloat(hashrate_ * rates["btc"])
	CNY := roundFloat(hashrate_ * rates["cny"])

	Commission := roundFloat(USD * config.DefaultPercentage)
	workerCommission := models.BillingWorkerMoney{Hashrate: hashrate, USD: USD, BTC: BTC, CNY: CNY, CommissionUSD: Commission, Worker: *worker}

	return b.BillingRepo.SaveWorkerMoney(workerCommission)
}

func roundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}
