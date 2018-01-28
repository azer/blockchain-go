package blockchain

import (
	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (bci *BlockchainIterator) Next() (*Block, error) {
	var block *Block

	err := bci.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.BlocksBucket))
		encodedBlock := bucket.Get(bci.CurrentHash)

		if deserialized, err := DeserializeBlock(encodedBlock); err == nil {
			block = deserialized
		} else {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	bci.CurrentHash = block.PrevBlockHash

	return block, nil
}
