package Services

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)

	if err != nil {
		fmt.Println("Error encoding block")
		panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	
	if err != nil {
		fmt.Println("Error decoding block")
		panic(err)
	}

	return &block
}