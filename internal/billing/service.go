package billing

import (
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/billing/utils"
	. "git.tor.ph/hiveon/pool/models"
	"github.com/mileusna/crontab"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type BillingCalculator struct {
	BillingRepo           *BillingRepository
}

func NewBillingCalculator() (*BillingCalculator) {
	return &BillingCalculator{BillingRepo:NewBillingRepository()}
}

func (b BillingCalculator) StartCalculation(er chan error) {
	ctab := crontab.New()
	err := ctab.AddJob("0 1 * * *", b.loadWorkerWalletStatistic)
	if err != nil {
		log.Error("Failed to start billing module: ", err)
		er <- err
		return
	}
}

func (b BillingCalculator)loadWorkerWalletStatistic() {
	WalletWorkerMapping := b.consumeMapping()
	b.consumeStatistic(WalletWorkerMapping)
}

func (b BillingCalculator) consumeMapping() map[string]string {
	log.Info("Consuming wallet/worker mapping started: ", time.Now())

	MappingAPI := config.MappingApi
	resp, err := http.Get(MappingAPI)
	if err != nil {
		log.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	WalletWorkerMapping :=  make(map[string]string)
	res := utils.ParseJSON(string(body), false)

	for k,v := range res {
		if k == "data" {
			value := v.([]interface{})
			for _,v1 := range value {
				value1 := v1.(map[string]interface{})
				worker:= value1["worker"].(string)
				wallet:= value1["wallet"].(string)
				if (worker != "" && wallet != "") {
					WalletWorkerMapping[worker] = wallet
				}
			}
		}
	}
	log.Info("Consuming wallet/worker mapping finished: ", time.Now())
	return WalletWorkerMapping
}

func (b BillingCalculator) consumeStatistic(WalletWorkerMapping map[string]string) {
	log.Info("Consuming wallet/worker statistic started: ", time.Now())

	WorkersAPI := config.WorkersAPI
	resp, err := http.Get(WorkersAPI)
	if err != nil {
		log.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	res := utils.ParseJSON(string(body), false)

	for k,v := range res {
		if k == "data" {
			value := v.([]interface{})
			for _,v1 := range value {
				value1 := v1.(map[string]interface{})

				worker:= value1["rig"].(string)
				validShares:= value1["validShares"].(float64)
				invalidShares:= value1["invalidShares"].(float64)
				staleShares:= value1["staleShares"].(float64)
				percentage:= value1["activityPercentage"].(float64)

				if (worker != "") {
					wallet := WalletWorkerMapping[worker]
					stat := BillingWorkerStatistic{ValidShares:validShares,InvalidShares:invalidShares,
						StaleShares:staleShares, ActivityPercentage:percentage}
					b.BillingRepo.SaveWorkerStatistic(stat, wallet, worker)
				}
			}
		}
	}
	log.Info("Consuming wallet/worker statistic finished: ", time.Now())
}

