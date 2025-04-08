package storage

import (
	"github.com/boltdb/bolt"
	"github.com/rootcontrol/blockchain/internal/domain"
)

type BoltRepository struct {
	db *bolt.DB
}

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
)

func NewBoltRepository() (*BoltRepository, error) {
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		return nil, err
	}

	repository := &BoltRepository{
		db: db,
	}

	return repository, nil
}

func (repo *BoltRepository) Close() {
	repo.db.Close()
}

func (repo *BoltRepository) SaveBlockchain(genesisBlock *domain.Block) error {
	err := repo.db.Update(func(tx *bolt.Tx) error {
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

func (repo *BoltRepository) InsertBlock(block *domain.Block) error {
	repo.db.Update(func(tx *bolt.Tx) error {
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

	return nil
}

func (repo *BoltRepository) GetLastHash() ([]byte, error) {
	var lastHash []byte

	err := repo.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			return nil
		}
		lastHash = bucket.Get([]byte("l"))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lastHash, nil
}

func (repo *BoltRepository) GetBlock(blockHash []byte) (*domain.Block, error) {
	var block *domain.Block

	err := repo.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(([]byte(blocksBucket)))
		encodedBlock := bucket.Get(blockHash)
		block = domain.DeserializeBLock(encodedBlock)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return block, nil
}