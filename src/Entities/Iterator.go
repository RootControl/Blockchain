package Entities

import (
	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Db          *bolt.DB
}