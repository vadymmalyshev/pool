package income

import (
	"git.tor.ph/hiveon/pool/internal/accounting"
	"time"

	"git.tor.ph/hiveon/pool/config"
	log "github.com/sirupsen/logrus"
)

type IncomeServicer interface {
	GetBlockCount() (BlockCount, error)
	GetIncomeHistory() (IncomeHistory, error)
}

type incomeService struct {
	accountingRepository accounting.AccointingRepositorer
	incomeRepository     IncomeRepositorer
}

func NewIncomeService() IncomeServicer {
	return &incomeService{accountingRepository: accounting.NewAccountingRepository(config.Seq2), incomeRepository: NewIncomeRepository(config.Seq3)}
}

func NewBlockServiceWithRepo(repo accounting.AccointingRepositorer) IncomeServicer {
	return &incomeService{accountingRepository: repo}
}

func (b *incomeService) GetBlockCount() (BlockCount, error) {
	var err error
	blockData := BlockCount{Code: 200}
	blockData.Data.Uncles, err = b.accountingRepository.GetUncleBlocks24h()
	if err != nil {
		return BlockCount{}, err
	}
	blockData.Data.Blocks, err = b.accountingRepository.GetNormalBlocks24h()
	if err != nil {
		return BlockCount{}, err
	}
	return blockData, nil
}

func (b *incomeService) GetIncomeHistory() (IncomeHistory, error) {
	rows, err := b.incomeRepository.GetIncomeHistory()
	if err != nil {
		return IncomeHistory{}, err
	}
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
	return incomeHistory, nil
}
