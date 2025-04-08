package services

import (
	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/domain"
)

type BlockchainService struct {
	Repository interfaces.BlockchainRepository
	Blockchain *domain.Blockchain
}

func NewBlockchainService(repo interfaces.BlockchainRepository) *BlockchainService {
	service := &BlockchainService{
		Repository: repo,
		Blockchain: domain.NewBlockchain(),
	}

	dbLastHash, err := service.Repository.GetLastHash()

	if err != nil || dbLastHash == nil {
		genesisBlock := domain.NewGenesisBlock()
		service.Repository.SaveBlockchain(genesisBlock)
		service.Blockchain.SetLastHash(genesisBlock.Hash)
	} else {
		service.Blockchain.SetLastHash(dbLastHash)
	}

	return service
}

func (service *BlockchainService) AddBlock(data string) error {
	block := domain.NewBlock(data, service.Blockchain.LastHash)

	err := service.Repository.InsertBlock(block)
	if err != nil {
		return err
	}

	service.Blockchain.SetLastHash(block.Hash)
	return nil
}
