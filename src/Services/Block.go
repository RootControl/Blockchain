
package Services

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp int64
	Data []byte
	PreviousBlockHash []byte
	Hash []byte
}

func (block *Block) CreateHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	
	headers := bytes.Join([][]byte {
		block.PreviousBlockHash,
		block.Data,
		timestamp,
	},
	[]byte{})

	hash := sha256.Sum256(headers)

	block.Hash = hash[:]
}

func NewBlock(data string, previousBlockHash []byte) *Block {
	block := &Block {
		time.Now().Unix(),
		[]byte(data),
		previousBlockHash,
		[]byte{},
	}

	block.CreateHash()

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}