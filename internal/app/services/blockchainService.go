package services

import (
	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/domain"
)

const genesisCoinbaseData = "In the grim darkness of the far future, the is only war!"

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

	coinbase := domain.NewCoinbaseTx(service.Address, genesisCoinbaseData)
	genesisBlock := domain.NewGenesisBlock(coinbase)
	service.Repository.SaveBlockchain(genesisBlock)
	bc.SetLastHash(genesisBlock.Hash)

	return bc
}

func (service *BlockchainService) AddBlock(data string) error {
	block := domain.NewBlock([]*domain.Transaction{}, service.Blockchain.LastHash)

	err := service.Repository.InsertBlock(block)
	if err != nil {
		return err
	}

	service.Blockchain.SetLastHash(block.Hash)
	return nil
}
