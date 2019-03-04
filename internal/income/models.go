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