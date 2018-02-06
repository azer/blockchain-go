package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().UnixNano(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (block *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(block); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (block *Block) HashTransactions() []byte {
	var (
		txHashes [][]byte
		txHash   [32]byte
	)

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.Id)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func DeserializeBlock(raw []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(raw))
	if err := decoder.Decode(&block); err != nil {
		return nil, err
	}

	return &block, nil
}
