package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"git.tor.ph/hiveon/pool/internal/consumer/model"
	. "git.tor.ph/hiveon/pool/internal/consumer/redis"
	. "git.tor.ph/hiveon/pool/internal/consumer/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var redisRepository IRedisRepository
var c *kafka.Consumer
var miningPools string

func StartConsumer(er chan error, buffer chan []byte, telebot *tb.Bot) {

	brokerList := GetConfig().GetString("kafka.brokers")
	caLocation := GetConfig().GetString("kafka.ca_location")
	user := GetConfig().GetString("kafka.username")
	password := GetConfig().GetString("kafka.password")
	groupId := GetConfig().GetString("kafka.group_id")
	topics := strings.Fields(GetConfig().GetString("kafka.topics"))
	miningPools = GetConfig().GetString("kafka.mining_pools")

	config := &kafka.ConfigMap{
		"api.version.request":  "true",
		"metadata.broker.list": brokerList,
		"security.protocol":    "sasl_ssl",
		"sasl.mechanisms":      "PLAIN",
		"ssl.ca.location":      caLocation,
		"sasl.username":        user,
		"sasl.password":        password,
		"group.id":             groupId,
		//"go.events.channel.enable":        true,
		//"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
	redisRepository = NewRedisRepository()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(config)

	if err != nil {
		log.Error("Failed to create consumer: ", err)
		er <- err
		os.Exit(1)
	}
	log.Info("Created consumer ", c)

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		log.Error("Can't subscribe to topics: " , topics)
		er <- err
		os.Exit(1)
	}
	log.Info("Subscribed to topics: ", topics)

	addTelebotKafkaEndpoints(telebot)

	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			fethDataFromKafka(c, buffer, er)
		}
	}
	log.Info("Closing consumer")

	c.Close()
}

func fethDataFromKafka(c *kafka.Consumer, buffer chan []byte, er chan error) {
	msg, err := c.ReadMessage(2000)
	if err == nil {
		proceedRow(string(msg.Value), buffer)
	}
}

func proceedRow(row string, buffer chan []byte) {
	res:= parseJSON(row)
	lastUpdateParams := map[string]interface{}{}
	delLastNotifyParams := []string{}

	messageType := res["msg_type"]
	messageMiner := res["miner"]
	var wallet, rig string
	if messageType == "payment-detail" {
		//TODO: payment details
	} else if messageMiner != nil{
		redisRepository.RedisCount("kafka") // save in redis
		//TODO: old format
	} else { // new format
		if res["serverNode"] != nil {
			if !checkNode(res["serverNode"].(string)) { // process the data from 3 nodes
				return
			}
		}
		redisRepository.RedisCount("kafka") // save in redis
		miner := createMiner(res, row)
		jsonRes,err := json.Marshal(miner)

		if err != nil {
			log.Error(err)
		}

		lastUpdateParams["0"] = miner.Tags["rig"]
		lastUpdateParams["1"], err = ParseTimestampToUnix(miner.Timestamp)
		if err!=nil {
			log.Error("Failed to parse the date: ", err)
		}

		rig = miner.Tags["rig"]
		wallet = miner.Tags["wallet"]

		buffer <- jsonRes
	}
	// 更新矿机最后提交时间
	if (len(lastUpdateParams) > 0) {
		redisRepository.RedisSet(wallet,lastUpdateParams)}

	if (len(delLastNotifyParams) > 0) {
		redisRepository.RedisDel(wallet, rig) }
}

func checkNode(s string) (bool){
	return strings.Contains(miningPools, s)
}

func parseJSON(row string) (map[string]interface{}) {
	res := make(map[string]interface{})
	f := map[string]interface{}{}

	err := json.Unmarshal([]byte(row), &f)
	if err != nil {
		log.Error(err)
	}

	for k, v := range f {
		switch v.(type) {
		case map[string]interface {}:
			m := v.(map[string]interface{})
			for k1, u := range m {
				if k == "sharesCount" {
					res["rig"] = k1
				}
				m1 := u.(map[string]interface{})
				for k2, u := range m1 {
					res[k2] = u
				}
			}
		default:
			res[k] = v
		}
	}
	return res
}

func createMiner(res map[string]interface{}, row string) (model.Miner){
	miner := model.Miner{};
	miner.Measurement = "worker"
	miner.Timestamp = res["createDt"].(string)
	Fields := make(map[string]interface{})
	Tags := make(map[string]string)

	hashrateCount := res["hashrateCount"]

	if hashrateCount != nil {
	if hashrateCount.(float64) < 1000 {
		Fields["LocalHashrate"] = hashrateCount.(float64)
	}}

	Fields["validShares"] = 0
	Fields["invalidShares"] = 0
	Fields["staleShares"] = 0
	Fields["originValidShares"] = 0

	if res["weightedCount"] != nil {
		Fields["validShares"] = res["weightedCount"].(float64)
	}

	if res["invalidCount"] != nil {
		Fields["invalidShares"] = res["invalidCount"].(float64)
	}

	if res["inferiorCount"] != nil {
		Fields["staleShares"] = res["inferiorCount"].(float64)
	}

	if res["validCount"] != nil {
		Fields["originValidShares"] = res["validCount"].(float64)
	}

	Tags["rig"] = ""
	Tags["token"] = ""
	Tags["wallet"] = ""

	if res["rig"] != nil {
		Tags["rig"] = res["rig"].(string)
	}

	if res["serverNode"] != nil {
		Tags["token"] = res["serverNode"].(string)
	}

	if res["minerWallet"] != nil {
		Tags["wallet"] = res["minerWallet"].(string)
	}

	miner.Fields = Fields
	miner.Tags = Tags

	return miner
}

func addTelebotKafkaEndpoints(telebot *tb.Bot) {
	telebot.Handle("/kafka", func(m *tb.Message) {
		message := IsUP
		/* _,er :=c.GetMetadata(nil, true,10000)

		if(er != nil) {
			message = IsDown
		}*/
		telebot.Send(m.Chat, "Kafka status : " + message)
	})
}

