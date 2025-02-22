package scheduler

import (
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/consumer/utils"
	"github.com/influxdata/influxdb1-client"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

var ethAPI, cnyAPI string
var retention, database, measurement, precision string
var clInflux *client.Client


func StartScheduler(er chan error, telebot *tb.Bot) {
	var err error

	ethAPI = config.Config.Scheduler.EthAPI
	cnyAPI = config.Config.Scheduler.CnyAPI

	retention = config.Config.Scheduler.Retention
	measurement = config.Config.Scheduler.Measurement
	database = config.Config.Kafka.DbName
	precision = config.Config.Kafka.Precision

	clInflux, err = config.Config.InfluxDB.Connect()

	if err != nil {
		log.Error(err)
		er <- fmt.Errorf("%s", err)
	}
	log.Info("Created influx client ", clInflux.Addr())

	fethCurrencyRates()
	c := cron.New()
	c.AddFunc("@every 7m", fethCurrencyRates)
	c.Start()
}

func fethCurrencyRates() {
	log.Info("Consuming the currency rates : ", time.Now())
	resp, err := http.Get(ethAPI)
	if err != nil {
		log.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	res := utils.ParseJSON(string(body), false)

	resp, err = http.Get(cnyAPI)
	if err != nil {
		log.Error(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	resCny := utils.ParseJSON(string(body), true)

	writeToInflux(res, resCny)
}

func writeToInflux(res map[string]interface{}, res2 map[string]interface{}) {
	var pts = make([]client.Point, 1)
	fields := make(map[string]interface{})
	fields["difficulty"] = res["difficulty"]
	fields["hashrate"] = res["hashrate"]
	fields["uncle_rate"] = res["uncle_rate"]

	cny := res2["price_cny"].(string)
	cnyFloat, _ := strconv.ParseFloat(cny, 64)
	cnyFloat = math.Round(cnyFloat*100) / 100 //round

	btc := res2["price_btc"].(string)
	btcFloat, _ := strconv.ParseFloat(btc, 64)

	fields["cny_float"] = cnyFloat
	fields["btc"] = btcFloat
	fields["usd"] = res["price_usd"]

	pts[0] = client.Point{
		Measurement: measurement,
		Fields:      fields,
		Time:        time.Now(),
	}

	bps := client.BatchPoints{
		Points:          pts,
		Database:        database,
		RetentionPolicy: retention,
		Precision:       precision,
		Time:            time.Now(),
	}
	_, err := clInflux.Write(bps)

	if err != nil {
		log.Error(err)
	}
}
