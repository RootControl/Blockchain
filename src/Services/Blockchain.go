package Services

import (
	"Blockchain/src/Entities"
	"Blockchain/src/Repositories"
)

const genesisCoinbaseData = "In the grim darkness of the far future, there is only war!"

func NewBlockchain(address string) *Entities.Blockchain {
	// Creates the first transaction in the blockchain
	coinbaseTransaction := NewCoinbaseTransaction(address, genesisCoinbaseData)

	genesisBlock := NewGenesisBlock(coinbaseTransaction)
	lastHash, db := Repositories.SaveGenesisBlockInDb(genesisBlock)

	return &Entities.Blockchain {
		LastHash: lastHash,
		Db: db,
	}
}

func AddDataToBlockchain(blockchain *Entities.Blockchain, data string) {
	lastHash := Repositories.GetLastHashFromDb(blockchain)
	transaction := NewCoinbaseTransaction("", data)

	newBlock := NewBlock([]*Entities.Transaction{transaction}, lastHash)
	
	Repositories.SaveBlockInDb(blockchain, newBlock)
}