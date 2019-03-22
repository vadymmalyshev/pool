package accounting

type RepoBillInfo struct {
	Balance   string
	FirstPaid float64
	FirstTime string
	TotalPaid float64
}

type Payment struct {
	firstPaid float64
	paymentId float64
}

func (r RepoBillInfo) isEmpty() bool {
	return r == RepoBillInfo{}
}
