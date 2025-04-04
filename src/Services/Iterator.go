package Services

import (
	"Blockchain/src/Entities"
	"Blockchain/src/Repositories"
)

func IterateNextBlock(iterator *Entities.BlockchainIterator) *Entities.Block {	
	var block *Entities.Block

	bockBytes := Repositories.GetBlockFromDb(iterator.Db, iterator.CurrentHash)
	block = DeserializeBlock(bockBytes)
	
	iterator.CurrentHash = block.PreviousBlockHash

	return block
}