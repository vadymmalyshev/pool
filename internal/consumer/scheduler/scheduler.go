package scheduler

import (
	"fmt"
	"github.com/robfig/cron"
	. "git.tor.ph/hiveon/pool/internal/consumer/utils"
	"github.com/influxdata/influxdb1-client"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ethAPI, cnyAPI string
var retention, database, measurement, precision string
var clInflux *client.Client

func StartScheduler(er chan error, telebot *tb.Bot) {
	ethAPI = GetConfig().GetString("scheduler.eth_API")
	cnyAPI = GetConfig().GetString("scheduler.cny_API")

	retention = GetConfig().GetString("scheduler.retention")
	measurement = GetConfig().GetString("scheduler.measurement")
	database = GetConfig().GetString("kafka.db_name")
	precision = GetConfig().GetString("kafka.precision")

	clInflux := getMinerdashClient()
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
	res := ParseJSON(string(body), false)

	resp, err = http.Get(cnyAPI)
	if err != nil {
		log.Error(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	resCny := ParseJSON(string(body), true)

	writeToInflux(res, resCny)
}

func getMinerdashClient() *client.Client {
	host := GetConfig().GetString("influx.host")
	port := GetConfig().GetString("influx.port")
	user := GetConfig().GetString("influx.username")
	password := GetConfig().GetString("influx.password")

	u, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))

	if err != nil {
		panic(err)
	}

	clInflux, err = client.NewClient(client.Config{URL: *u})
	if err != nil {
		log.Error(err)
	}

	if _, _, err := clInflux.Ping(); err != nil {
		log.Error(err)
	}

	clInflux.SetAuth(user, password)

	return clInflux
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
	//cnyString := strconv.FormatFloat(cnyFloat, 'f', 2, 64)

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
