package blockchain

import (
	"github.com/boltdb/bolt"
)

var DB *bolt.DB

func init() {
	db, err := bolt.Open(config.DBFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	DB = db
}

func LastHash() ([]byte, error) {
	var lastHash []byte

	err := DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.BlocksBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})

	return lastHash, err
}

func Tip() ([]byte, error) {
	var tip []byte

	err := DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.BlocksBucket))

		if bucket != nil {
			tip = bucket.Get([]byte("l"))
			return nil
		}

		genesis := NewGenesisBlock()
		bucket, err := tx.CreateBucket([]byte(config.BlocksBucket))
		if err != nil {
			return err
		}

		serialized, err := genesis.Serialize()
		if err != nil {
			return err
		}

		if err := bucket.Put(genesis.Hash, serialized); err != nil {
			return err
		}

		if err := bucket.Put([]byte("l"), genesis.Hash); err != nil {
			return err
		}

		tip = genesis.Hash

		return nil
	})

	return tip, err
}
