package interfaces

import "github.com/rootcontrol/blockchain/internal/domain"

type BlockchainRepository interface {
	SaveBlockchain(genesisBlock *domain.Block) error
	InsertBlock(block *domain.Block) error
	GetLastHash() ([]byte, error)
	GetBlock(blockHash []byte) (*domain.Block, error)
	Close()
	SaveWallet(wallet *domain.Wallet) error
	GetWallet(address string) (*domain.Wallet, error)
}