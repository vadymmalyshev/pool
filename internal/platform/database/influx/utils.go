package influx

import "encoding/json"

func ToFloat64(val interface{}) float64 {
	result, _ := val.(json.Number).Float64()
	return result
}
