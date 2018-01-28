package blockchain

import (
	"github.com/boltdb/bolt"
)

func NewBlockChain() (*Blockchain, error) {
	tip, err := Tip()
	if err != nil {
		return nil, err
	}

	return &Blockchain{
		Tip: tip,
		DB:  DB,
	}, nil
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) error {
	lastHash, err := LastHash()
	if err != nil {
		return err
	}

	newBlock := NewBlock(data, lastHash)

	return bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.BlocksBucket))

		serialized, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		if err := bucket.Put(newBlock.Hash, serialized); err != nil {
			return err
		}

		if err := bucket.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}

		bc.Tip = newBlock.Hash

		return nil
	})
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: bc.Tip,
		DB:          bc.DB,
	}
}
