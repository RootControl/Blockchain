package Services

import (
	"bytes"
	"encoding/gob"
	"time"
	"Blockchain/src/Entities"
)

func NewGenesisBlock() *Entities.Block {
	return NewBlock("Genesis Block", []byte{})
}

// TODO: remove previousBlockHash from the function
func NewBlock(data string, previousBlockHash []byte) *Entities.Block {
	block := &Entities.Block {
		Timestamp: time.Now().Unix(),
		Data:[]byte(data),
		PreviousBlockHash: previousBlockHash,
		Hash: []byte{},
		Nonce:0,
	}

	pow := CreateProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func DeserializeBlock(data []byte) *Entities.Block {
	var block Entities.Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	
	if err != nil {
		panic(err)
	}

	return &block
}