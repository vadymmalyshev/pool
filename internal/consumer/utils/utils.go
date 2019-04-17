package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

var (
	v *viper.Viper
)

func init() {
	v = viper.New()
	v.SetConfigType("yaml")

	v.AddConfigPath(".././conf/")
	v.AddConfigPath("./conf/")

	err := v.ReadInConfig()
	if err != nil {
		log.Error(err)
	}
}

func GetConfig() *viper.Viper {
	return v
}

func ParseJSON(row string, replace bool) (map[string]interface{}) {
	res := make(map[string]interface{})
	f := map[string]interface{}{}

	if (replace) {
		row = strings.Replace(row, "[","", -1)
		row = strings.Replace(row, "]","", -1)
	}

	err := json.Unmarshal([]byte(row), &f)
	if err != nil {
		log.Error(err)
	}

	for k, v := range f {
		switch v.(type) {
		case map[string]interface {}:
			m := v.(map[string]interface{})
			for k1, u := range m {
					res[k1] = u
			}
		default:
			res[k] = v
		}
	}
	return res
}

func ParseTimestampToUnix(stringTime string) (int64, error) {
	splitted_string := strings.Split(stringTime,".")
	right_part := ""

	if (len(splitted_string)==2){
		right_part = splitted_string[1]
	}
	left_part := splitted_string[0]
	right_part_int, _ := strconv.Atoi(right_part)

	fz := time.FixedZone("CST", 8*3600) // China time
	timestamp,error := time.ParseInLocation("2006-01-02T15:04:05", left_part, fz)
	tz,_ := time.LoadLocation("UTC")
	timestamp_res := timestamp.In(tz).UTC().UnixNano()

	//delay4m:= time.Minute * 4 // some delay from origin consumer, need to clarify
	timestamp_res = timestamp_res + int64(right_part_int) //+ int64(delay4m)

	return timestamp_res, error
}