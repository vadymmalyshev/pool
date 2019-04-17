package influx

import (
	"encoding/json"
	"fmt"
	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/consumer/model"
	"git.tor.ph/hiveon/pool/internal/consumer/redis"
	"git.tor.ph/hiveon/pool/internal/consumer/utils"
	"github.com/influxdata/influxdb1-client"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var clInflux *client.Client
var redisRepository redis.IRedisRepository

var retention string
var database string
var precision string


func StartDBProducer(er chan error, buffer chan []byte, telebot *tb.Bot) {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	redisRepository = redis.NewRedisRepository()
	//addTelebotRedisEndpoints(telebot)

	retention = utils.GetConfig().GetString("kafka.retention")
	database = utils.GetConfig().GetString("kafka.db_name")
	precision = utils.GetConfig().GetString("kafka.precision")

	clInflux := getMinerdashClient();
	//addTelebotInfluxEndpoints(telebot)

	log.Info("Created influx client ", clInflux.Addr())

	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		case data := <-buffer:
			writeToInflux(data)
		}
	}
}

func getMinerdashClient() (*client.Client) {

	host := config.InfluxDB.Host
	port := config.InfluxDB.Port
	user := config.InfluxDB.User
	password := config.InfluxDB.Pass

	u, err := url.Parse(fmt.Sprintf("http://%s:%s", host, strconv.Itoa(port)))

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

func writeToInflux(data []byte) {
	var pts = make([]client.Point, 1)

	res := model.Miner{}
	err := json.Unmarshal([]byte(data), &res)
	if err!=nil {
		log.Error("Failed to unmarshall the data: ", err)
	}

	stringTime := res.Timestamp
	splittedString := strings.Split(stringTime,".")
	rightPart := ""

	if (len(splittedString)==2){
		rightPart = splittedString[1]
	}

	leftPart := splittedString[0]
	rightPartInt, _ := strconv.Atoi(rightPart)

   	fz := time.FixedZone("CST", 8*3600) // China time
	timestamp,error := time.ParseInLocation("2006-01-02T15:04:05", leftPart, fz)
	timestamp.Add(time.Minute * 4) //DELAY_MINUTES
	timestamp.Add(time.Nanosecond * time.Duration(rightPartInt))
	tz,err := time.LoadLocation("UTC")
	timestamp = timestamp.In(tz)

	if error!=nil {
		log.Error("Failed to parse the date: ", error)
	}

	pts[0] = client.Point{
		Measurement: res.Measurement,
		Tags: res.Tags,
		Fields: res.Fields,
		Time: timestamp,
	}

	bps := client.BatchPoints{
		Points:          pts,
		Database:        database,
		RetentionPolicy: retention,
		Precision: precision,
		Time: time.Now(),
	}
	_, err = clInflux.Write(bps)

	if err != nil {
		log.Error(err)
	}
	redisRepository.RedisCount("w_influx");
	log.Info("Added a point to Influx : ", time.Now());
}

func addTelebotRedisEndpoints(telebot *tb.Bot) {
	telebot.Handle("/redis", func(m *tb.Message) {
		telebot.Send(m.Chat, "Redis status : " + redisRepository.RedisAlive())
	})

	telebot.Handle("/redis-stat", func(m *tb.Message) {
		now := time.Now()
		k := []string {strconv.Itoa(int(now.Month())),"-",strconv.Itoa(now.Day())}
		key := strings.Join(k,"")
		keyArray :=[]string {"count:kafka", ":", key}
		hashId := strings.Join(keyArray,"")

		telebot.Send(m.Chat, "Redis statistics by minutes  : ")
		res := redisRepository.RedisGetPoints(hashId)
		for k,v := range res {
			telebot.Send(m.Chat, k + " " +v)
		}
	})
}

func addTelebotInfluxEndpoints(telebot *tb.Bot) {
	telebot.Handle("/influx", func(m *tb.Message) {
		message := utils.IsUP
		_,_,er := clInflux.Ping()
		if(er != nil) {
			message = utils.IsDown
		}
		telebot.Send(m.Chat, "Influx status : " + message)
	})
}