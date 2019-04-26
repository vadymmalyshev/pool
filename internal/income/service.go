package income

import (
	"github.com/jinzhu/gorm"
	"time"

	"git.tor.ph/hiveon/pool/internal/accounting"
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

func NewIncomeService(sql2DB *gorm.DB, sql3DB *gorm.DB) IncomeServicer {
	return &incomeService{accountingRepository: accounting.NewAccountingRepository(sql2DB), incomeRepository: NewIncomeRepository(sql3DB)}
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
