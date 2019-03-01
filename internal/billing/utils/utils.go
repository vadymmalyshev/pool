package utils

import (
	"strings"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

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
		case []interface{}:
			m := v.([]interface{})
			res[k] = m
		default:
			res[k] = v
		}
	}
	return res
}
