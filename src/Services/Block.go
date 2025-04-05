package Services

import (
	"bytes"
	"encoding/gob"
	"time"
	"Blockchain/src/Entities"
)

func NewGenesisBlock(coinbase *Entities.Transaction) *Entities.Block {
	return NewBlock([]*Entities.Transaction{coinbase}, []byte{})
}

// TODO: remove previousBlockHash from the function
func NewBlock(transactions []*Entities.Transaction, previousBlockHash []byte) *Entities.Block {
	block := &Entities.Block {
		Timestamp: time.Now().Unix(),
		Transactions: transactions,
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
	if len(data) == 0 {
		return nil
	}

	var block Entities.Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	
	if err != nil {
		panic(err)
	}

	return &block
}