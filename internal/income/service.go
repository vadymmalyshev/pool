package income

import (
	. "git.tor.ph/hiveon/pool/internal/accounting"
	"github.com/jinzhu/gorm"
)

type BlockServicer interface {
	GetBlockCount() BlockCount
}

type blockService struct {
	accountingRepository AccointingRepositorer
}

func NewBlockService(db *gorm.DB) BlockServicer {
	return &blockService{accountingRepository: NewAccountingRepository(db)}
}

func NewBlockServiceWithRepo(repo AccointingRepositorer) BlockServicer {
	return &blockService{accountingRepository: repo}
}

func (b *blockService) GetBlockCount() BlockCount {
	blockData := BlockCount{Code: 200}
	blockData.Data.Uncles = b.accountingRepository.GetBlock24Uncle()
	blockData.Data.Blocks = b.accountingRepository.GetBlock24NotUnckle()
	return blockData
}
