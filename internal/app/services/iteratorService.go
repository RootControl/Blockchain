package services

import (
	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/domain"
)

type IteratorService struct {
	Repository interfaces.BlockchainRepository
	Iterator   *domain.BlockchainIterator
}

func NewIteratorService(repository interfaces.BlockchainRepository, blockchainLastHash []byte) *IteratorService {
	return &IteratorService{
		Repository: repository, 
		Iterator: domain.NewBlockchainIterator(blockchainLastHash),
	}
}

func (service *IteratorService) NextBlock() *domain.Block {
	block, err := service.Repository.GetBlock(service.Iterator.CurrentHash)

	if err != nil {
		return nil
	}

	service.Iterator.CurrentHash = block.PrevBlockHash

	return block
}