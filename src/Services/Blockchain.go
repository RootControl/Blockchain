package Services

import (
	"Blockchain/src/Entities"
)

func NewBlockchain() *Entities.Blockchain {
	return &Entities.Blockchain{
		Blocks: []*Entities.Block{NewGenesisBlock()},
	}
}