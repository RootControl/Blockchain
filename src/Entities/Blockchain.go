package Entities

import (
	"github.com/boltdb/bolt"
)

type Blockchain struct {
	LastHash []byte
	Db     *bolt.DB
}