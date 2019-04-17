package billing

import (
	"fmt"
	"time"
	"git.tor.ph/hiveon/pool/models"
)

const timeLayout = "02.01.2006"

// Collector collects workers fee.
type Collector struct {
	billingRepository *BillingRepository
	calculator        *BillingCalculator
}

// NewCollector returns collector.
func NewCollector(billingRepository *BillingRepository) *Collector {
	return &Collector{billingRepository: billingRepository}
}

// ChargeWorkers charges worker and save into DB.
func (c Collector) ChargeWorkers(charge models.Charge) error {
	workers := make([]models.WorkerFee, len(charge.WorkersFee))
	for k, v := range charge.WorkersFee {
		v.Date = charge.Date.Format(timeLayout)
		v.Coin = models.Coin{Name: v.CoinName}
		workers[k] = v
	}

	return c.billingRepository.BulkUpdateWorkersFeeIfNotExist(workers)
}

// CalculateFees charges worker and save into DB.
func (c Collector) CalculateFees(date time.Time) ([]models.Bill, error) {

	fees, err := c.billingRepository.FindWorkersFeeByDate(date.Format(timeLayout))
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(fees))
	for k, v := range fees {
		ids[k] = v.WorkerID
	}

	stats, err := c.billingRepository.AllWorkerStatistic(ids, date.Format(timeLayout))
	if err != nil {
		return nil, err
	}

	if len(stats) == 0 {
		return nil, nil
	}

	billList := c.billList(stats)
	for _, bl := range billList {
		for k := range fees {
			if bl.WalletAdds == fees[k].WalletAddr {
				fees[k].Shares = bl.Amount
			}
		}
	}

	if err := c.billingRepository.BulkUpdateWorkersShares(fees); err != nil {
		return nil, err
	}

	return billList, nil
}

func (c Collector) billList(stats []models.BillingWorkerStatistic) []models.Bill {
	bills := make(map[string]models.Bill)
	for _, v := range stats {
		if _, ok := bills[v.Wallet.Address]; !ok {
			bills[v.Wallet.Address] = models.Bill{
				WalletAdds: v.Wallet.Address,
				Coin:       v.Wallet.Coin.Name,
			}
		}
		sl := bills[v.Wallet.Address]
		sl.Workers = append(sl.Workers, models.BillWorker{
			ID: fmt.Sprintf("%s#id%d",
				v.Worker.Name,
				v.WorkerID,
			),
			Shares: v.ValidShares})
		bills[v.Wallet.Address] = sl
	}

	billList := make([]models.Bill, 0, len(bills))
	for _, b := range bills {
		var amount float64
		for _, v := range b.Workers {
			amount += v.Shares
		}
		b.Amount = amount
		billList = append(billList, b)
	}

	return billList
}
