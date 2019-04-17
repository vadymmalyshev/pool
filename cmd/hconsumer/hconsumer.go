package main

import (
	"fmt"
	"git.tor.ph/hiveon/pool/internal/consumer/influx"
	"git.tor.ph/hiveon/pool/internal/consumer/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq"
	"git.tor.ph/hiveon/pool/internal/consumer/scheduler"
	"git.tor.ph/hiveon/pool/internal/consumer/telebot"
	"git.tor.ph/hiveon/pool/internal/consumer/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	var f *os.File
	var err error

	errs := make(chan error, 2)
	buffer := make(chan []byte, 10)

	// logrus
	if f, err = os.Create("testlogrus" + ".log"); err != nil {
		log.Error("creat log failed", err)
		return
	}

	defer func(f *os.File) {
		f.Sync()
		f.Close()
	}(f)

	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "02-01-2006 15:04:05", PrettyPrint: true})
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw) // multilogging, both stdout and file

	conf := utils.GetConfig()

	telebot.CreateBot(
		conf.GetString("telegrambot.token"),
		conf.GetInt64("telegrambot.chatID"))

	createLoggerEnpoint(telebot.Bot)

	go kafka.StartConsumer(errs, buffer, telebot.Bot)
	go influx.StartDBProducer(errs, buffer, telebot.Bot)
	go scheduler.StartScheduler(errs, telebot.Bot)
	go telebot.StartBot()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Info("terminated", <-errs)
}

func createLoggerEnpoint(telebot *tb.Bot) {

	telebot.Handle("/stat-1h", func(m *tb.Message) {

		fileData, _ := ioutil.ReadFile("testlogrus.log")
		res:= proceedJsonData(fileData)

		message := "no data"
		startTime := time.Now().Add(time.Hour *-1 )

		count := 0
		for _,v:= range res {
			recordTime :=v.(map[string]interface{})["time"].(string)
			formattedTime, _:=time.Parse("02-01-2006 15:04:05", recordTime)
			if inTimeSpan(startTime, formattedTime) {
				count ++
			}
		}
		if count > 0 {
			message = strconv.Itoa(count)
		}

		telebot.Send(m.Chat, "Added data points number in Influx during 1 hour: "+message)
	})
}

func proceedJsonData(fileData []byte ) ([]interface{}){
	jsonArray :=[]string {"[", bytesToString(fileData),"]"}
	str := strings.Join(jsonArray," ")
	str = strings.Replace(str, "}\n{","},\n{",-1)
	jq := gojsonq.New().JSONString(str)
	res := jq.WhereContains("msg", "Added a point to Influx").Select("time").Get()
	r := res.([]interface{})
	return r
}

func inTimeSpan(start, check time.Time) bool {
	return check.After(start)
}

func bytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

