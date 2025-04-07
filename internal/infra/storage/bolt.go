package storage

import (
	"github.com/boltdb/bolt"
	"github.com/rootcontrol/blockchain/internal/domain"
)

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
)

func SaveBlockchain(genesisBlock *domain.Block) error {
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			bucket, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}

			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				return err
			}

			err = bucket.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func InsertBlock(block *domain.Block) error {
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		return err
	}

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			return err
		}

		err = bucket.Put([]byte("l"), block.Hash)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func GetLastHash() ([]byte, error) {
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		return nil, err
	}

	var lastHash []byte

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(([]byte(blocksBucket)))
		lastHash = bucket.Get([]byte("l"))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lastHash, nil
}