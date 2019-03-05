package income

import (
	"git.tor.ph/hiveon/pool/config"
	. "git.tor.ph/hiveon/pool/internal/accounting"
	log "github.com/sirupsen/logrus"
	"time"
)

type IncomeServicer interface {
	GetBlockCount() BlockCount
	GetIncomeHistory() IncomeHistory
}

type incomeService struct {
	accountingRepository AccointingRepositorer
	incomeRepository     IncomeRepositorer
}

func NewIncomeService() IncomeServicer {
	return &incomeService{accountingRepository: NewAccountingRepository(config.Seq2), incomeRepository: NewIncomeRepository(config.Seq3)}
}

func NewBlockServiceWithRepo(repo AccointingRepositorer) IncomeServicer {
	return &incomeService{accountingRepository: repo}
}

func (b *incomeService) GetBlockCount() BlockCount {
	blockData := BlockCount{Code: 200}
	blockData.Data.Uncles = b.accountingRepository.GetBlock24Uncle()
	blockData.Data.Blocks = b.accountingRepository.GetBlock24NotUnckle()
	return blockData
}

func (b *incomeService) GetIncomeHistory() IncomeHistory {
	rows := b.incomeRepository.GetIncomeHistory()
	var incomeSlice []Income

	for rows.Next() {
		var income Income
		var t int64
		err := rows.Scan(&t, &income.Income)
		if err != nil {
			log.Error(err)
		}
		income.Time = time.Unix(t, 0).Format(time.RFC3339)
		incomeSlice = append(incomeSlice, income)
	}

	incomeHistory := IncomeHistory{Code: 200}
	incomeHistory.Data = incomeSlice
	return incomeHistory
}

