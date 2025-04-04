package Entities

import (
	"github.com/boltdb/bolt"
)

type Blockchain struct {
	LastHash []byte
	Db     *bolt.DB
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator {
		CurrentHash: blockchain.LastHash,
		Db: blockchain.Db,
	}
}