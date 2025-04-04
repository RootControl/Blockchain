package Repositories

import (
	"Blockchain/src/Entities"
	"github.com/boltdb/bolt"
)

const blockBucket = "blocks"
const dbFile = "blockchain.db"

func SaveGenesisBlockInDb(genesisBlock *Entities.Block) ([]byte, *bolt.DB) {
	var lastHash []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		panic(err)
	}

	db.Update(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(blockBucket))

			if bucket != nil {
				lastHash = bucket.Get([]byte("l"))
			} else {
				bucket, err := tx.CreateBucket([]byte(blockBucket))
				if err != nil {
					panic(err)
				}

				bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
				bucket.Put([]byte("l"), genesisBlock.Hash)

				lastHash = genesisBlock.Hash
			}

			return nil
		},
	)

	return lastHash, db
}

func SaveBlockInDb(blockchain *Entities.Blockchain, block *Entities.Block) {
	blockchain.Db.Update(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(blockBucket))
			bucket.Put(block.Hash, block.Serialize())
			bucket.Put([]byte("l"), block.Hash)
			blockchain.LastHash = block.Hash
			return nil
		},
	)
}

func GetLastHashFromDb(blockchain *Entities.Blockchain) []byte {
	var lastHash []byte

	blockchain.Db.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(blockBucket))
			lastHash = bucket.Get([]byte("l"))

			return nil
		},
	)

	return lastHash
}

func GetBlockFromDb(database *bolt.DB, hash []byte) []byte {
	var block []byte

	database.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(blockBucket))
			block = bucket.Get(hash)

			return nil
		},
	)

	return block
}
