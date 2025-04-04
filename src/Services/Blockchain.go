package Services

import (
	"Blockchain/src/Entities"
	"Blockchain/src/Repositories"
)

func NewBlockchain() *Entities.Blockchain {
	genesisBlock := NewGenesisBlock()
	lastHash, db := Repositories.SaveGenesisBlockInDb(genesisBlock)

	return &Entities.Blockchain {
		LastHash: lastHash,
		Db: db,
	}
}

func AddDataToBlockchain(blockchain *Entities.Blockchain, data string) {
	lastHash := Repositories.GetLastHashFromDb(blockchain)
	
	newBlock := NewBlock(data, lastHash)
	
	Repositories.SaveBlockInDb(blockchain, newBlock)
}