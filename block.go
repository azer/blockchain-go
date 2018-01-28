package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"
)

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().UnixNano(),
		Data:          []byte(data),
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
	Data          []byte
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

func DeserializeBlock(raw []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(raw))
	if err := decoder.Decode(&block); err != nil {
		return nil, err
	}

	return &block, nil
}
