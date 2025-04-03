
package Services

import (
	"time"
)

type Block struct {
	Timestamp int64
	Data []byte
	PreviousBlockHash []byte
	Hash []byte
	Nonce int
}

func NewBlock(data string, previousBlockHash []byte) *Block {
	block := &Block {
		time.Now().Unix(),
		[]byte(data),
		previousBlockHash,
		[]byte{},
		0,
	}

	pow := CreateProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}