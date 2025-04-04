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

func GetAllBlocksFromDb(blockchain *Entities.Blockchain) [][]byte {
	var blocks [][]byte

	blockchain.Db.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(blockBucket))

			bucket.ForEach(
				func(key, value []byte) error {
					// Skip the "l" key which stores the last hash
					if string(key) != "l" {
						blocks = append(blocks, value)
					}
					return nil
				},
			)

			return nil
		},
	)

	return blocks
}