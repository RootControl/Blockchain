package Entities

type Block struct {
	Timestamp         int64
	Data              []byte
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
}
