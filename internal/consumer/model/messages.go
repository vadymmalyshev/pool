package model

type Miner struct {
	Measurement string `json:"measurement"`
	Tags map[string] string `json:"tags"`
	Fields map[string] interface{} `json:"fields"`
	Timestamp string `json:"time"`
}