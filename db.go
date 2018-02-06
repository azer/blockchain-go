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

func AlreadyCreated() bool {
	created := false

	err := DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.BlocksBucket))
		created = bucket != nil
		return nil
	})

	return err == nil && created
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

func GetTip() ([]byte, error) {
	var tip []byte

	err := DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(config.BlocksBucket))
		if err != nil {
			return err
		}

		tip = bucket.Get([]byte("l"))

		return nil
	})

	return tip, err
}

func CreateTipIfNotExists(genesisBlock *Block) ([]byte, error) {
	var tip []byte

	err := DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.BlocksBucket))

		if bucket != nil {
			tip = bucket.Get([]byte("l"))
			return nil
		}

		bucket, err := tx.CreateBucket([]byte(config.BlocksBucket))
		if err != nil {
			return err
		}

		serialized, err := genesisBlock.Serialize()
		if err != nil {
			return err
		}

		if err := bucket.Put(genesisBlock.Hash, serialized); err != nil {
			return err
		}

		if err := bucket.Put([]byte("l"), genesisBlock.Hash); err != nil {
			return err
		}

		tip = genesisBlock.Hash

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tip, err
}
