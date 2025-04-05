package Entities

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Block struct {
	Timestamp         int64
	Transactions      []*Transaction
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)

	if err != nil {
		panic(err)
	}

	return result.Bytes()
}

func (block *Block) HashTransactions() []byte {
	var transactionHashes [][]byte

	for _, transaction := range block.Transactions {
		transactionHashes = append(transactionHashes, transaction.Id)
	}

	transactionHash := sha256.Sum256(bytes.Join(transactionHashes, []byte{}))

	return transactionHash[:]
}