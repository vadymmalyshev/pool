package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/influxdata/influxdb1-client/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var (
	v *viper.Viper
)

const DefaultTimeFormat = "2006-01-02 15:04:05"

func init() {
	v = viper.New()
	v.SetConfigType("yaml")

	v.AddConfigPath(".././config/")
	v.AddConfigPath("./config/")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	v.SetEnvPrefix(v.GetString("appname"))
	if err != nil {
		log.Fatal(err)
	}
}
func GetConfig() *viper.Viper {
	return v
}

func FormatTimeToRFC3339(inputTime string) string {
	rfcTime, _ := time.Parse(DefaultTimeFormat, inputTime)
	return rfcTime.Format(time.RFC3339)
}

func GetRowStringValue(row Row, index int, columnName string) string {
	return GetRowValue(row, index, columnName).(string)
}

func GetRowFloatValue(row Row, index int, columnName string) float64 {
	res, err := GetRowValue(row, index, columnName).(json.Number).Float64()
	if err != nil {
		return 0
	}
	return res
}

//return zero in case of any inconsistencies
func GetRowValue(row Row, index int, columnName string) interface{} {
	if row.Values[index] == nil {
		return 0
	}

	colIndex, err := getIndexByColumnName(row, columnName)
	if err != nil {
		log.Error(err)
	}

	if row.Values[index][colIndex] != nil {
		res := row.Values[index][colIndex]
		return res
	}
	return json.Number(0)
}

func getIndexByColumnName(row Row, columnName string) (int, error) {

	for i, name := range row.Columns {
		if name == columnName {
			return i, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Column with name %s hasn't been found", columnName))
}


func RoundFloat2(value float64) float64{
	return float64(int(value *100))/100
}

func FormatWalletID(wallet string) string {
	s := strings.ToLower(wallet)
	if (strings.HasPrefix(s, "0xwallet")) {
		strings.Replace(wallet,"0xwallet","wallet",1)
	}
	return s
}

func FormatWorkerName(name string) string {
	splitted_string := strings.Split(name,"#id")
	return splitted_string[0]
}
