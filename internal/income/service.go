package income

import (
	. "git.tor.ph/hiveon/pool/internal/api/repository"
	"git.tor.ph/hiveon/pool/internal/api/response"
)

type BlockServicer interface {
	GetBlockCount() response.BlockCount
}

type blockService struct {
	hiveosRepository IHiveosRepository
}

func NewBlockService() BlockService {
	return &blockService{hiveosRepository: NewHiveosRepository()}
}

func NewBlockServiceWithRepo(repo IHiveosRepository) BlockService {
	return &blockService{hiveosRepository: repo}
}

func (b *blockService) GetBlockCount() response.BlockCount {
	blockData := response.BlockCount{Code: 200}
	blockData.Data.Uncles = b.hiveosRepository.GetBlock24Uncle()
	blockData.Data.Blocks = b.hiveosRepository.GetBlock24NotUnckle()
	return blockData
}
