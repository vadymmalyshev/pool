package income

import (
	. "git.tor.ph/hiveon/pool/internal/accounting"
	"git.tor.ph/hiveon/pool/internal/api/response"
)

type BlockServicer interface {
	GetBlockCount() response.BlockCount
}

type blockService struct {
	accountingRepository AccointingRepositorer
}

func NewBlockService() BlockServicer {
	return &blockService{accountingRepository: NewAccountingRepository()}
}

func NewBlockServiceWithRepo(repo AccointingRepositorer) BlockServicer {
	return &blockService{accountingRepository: repo}
}

func (b *blockService) GetBlockCount() response.BlockCount {
	blockData := response.BlockCount{Code: 200}
	blockData.Data.Uncles = b.accountingRepository.GetBlock24Uncle()
	blockData.Data.Blocks = b.accountingRepository.GetBlock24NotUnckle()
	return blockData
}
