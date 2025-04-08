package services

import (
	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/domain"
)

type IteratorService struct {
	repository interfaces.BlockchainRepository
	Iterator   *domain.BlockchainIterator
}

func NewIteratorService(repository interfaces.BlockchainRepository, iterator domain.BlockchainIterator) *IteratorService {
	return &IteratorService{repository: repository, Iterator: &iterator}
}

func (service *IteratorService) NextBlock() *domain.Block {
	block, err := service.repository.GetBlock(service.Iterator.CurrentHash)

	if err != nil {
		return nil
	}

	service.Iterator.CurrentHash = block.PrevBlockHash

	return block
}