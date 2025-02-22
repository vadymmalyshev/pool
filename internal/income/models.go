package income

// Success response
// swagger:response BlockCount
type BlockCount struct {
	Code int `json:"code"`
	Data struct {
		Uncles int `json:"uncles"`
		Blocks int `json:"blocks"`
	} `json:"data"`
}

// Success response
// swagger:response IncomeHistory
type IncomeHistory struct {
	Code int      `json:"code"`
	Data []Income `json:"data"`
}

type Income struct {
	Time   string `json:"time"`
	Income string `json:"income"`
}