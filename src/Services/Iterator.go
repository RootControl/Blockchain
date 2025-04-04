package Services

import (
	"Blockchain/src/Entities"
	"Blockchain/src/Repositories"
)

func IterateNextBlock(iterator *Entities.BlockchainIterator) *Entities.Block {	
	blockBytes := Repositories.GetBlockFromDb(iterator.Db, iterator.CurrentHash)
	
	if blockBytes == nil {
		return nil
	}
	
	block := DeserializeBlock(blockBytes)
	iterator.CurrentHash = block.PreviousBlockHash
	return block
}