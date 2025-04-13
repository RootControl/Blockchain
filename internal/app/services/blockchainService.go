package services

import (
	"errors"

	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/domain"
)

type BlockchainService struct {
	Repository interfaces.BlockchainRepository
	Blockchain *domain.Blockchain
	Address    string
}

func NewBlockchainService(repo interfaces.BlockchainRepository, address string) *BlockchainService {
	service := &BlockchainService{
		Repository: repo,
		Blockchain: domain.NewBlockchain(),
		Address:    address,
	}

	dbLastHash, err := service.Repository.GetLastHash()

	if err != nil || dbLastHash == nil {
		service.Blockchain = service.CreateBlockchain()
	} else {
		service.Blockchain.SetLastHash(dbLastHash)
	}

	return service
}

func (service *BlockchainService) CreateBlockchain() *domain.Blockchain {
	bc := domain.NewBlockchain()

	coinbase := domain.NewCoinbaseTx(service.Address)
	genesisBlock := domain.NewGenesisBlock(coinbase)
	service.Repository.SaveBlockchain(genesisBlock)
	bc.SetLastHash(genesisBlock.Hash)

	return bc
}

func (service *BlockchainService) MineBlock(transactions []*domain.Transaction) error {
	txService := NewTransactionService(service.Repository, service.Blockchain.LastHash)

	for _, tx := range transactions {
		if !txService.VerifyTransaction(tx) {
			return errors.New("ERROR: Invalid transaction")
		}
	}

	block := domain.NewBlock(transactions, service.Blockchain.LastHash)

	err := service.Repository.InsertBlock(block)
	if err != nil {
		return err
	}

	service.Blockchain.SetLastHash(block.Hash)
	return nil
}
