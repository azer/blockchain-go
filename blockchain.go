package blockchain

import (
	"encoding/hex"
	"errors"
	"github.com/boltdb/bolt"
)

const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func NewBlockchain(rewardAddress string) (*Blockchain, error) {
	if AlreadyCreated() {
		return nil, errors.New("Already created before")
	}

	cbtx, err := NewCoinbaseTX(rewardAddress, genesisCoinbaseData)
	if err != nil {
		return nil, err
	}

	genesis := NewGenesisBlock(cbtx)

	tip, err := CreateTipIfNotExists(genesis)
	if err != nil {
		return nil, err
	}

	return &Blockchain{
		Tip: tip,
		DB:  DB,
	}, nil
}

func RestoreBlockchain() (*Blockchain, error) {
	if !AlreadyCreated() {
		return nil, errors.New("Database hasn't created yet")
	}

	tip, err := GetTip()
	if err != nil {
		return nil, err
	}

	return &Blockchain{
		Tip: tip,
		DB:  DB,
	}, nil

}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) error {
	lastHash, err := LastHash()
	if err != nil {
		return err
	}

	newBlock := NewBlock(transactions, lastHash)

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

func (blockchain *Blockchain) FindUnspentTransactions(address string) ([]*Transaction, error) {
	var unspentTXs []*Transaction
	spentTXOs := make(map[string][]int)
	bci := blockchain.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			encodedTxId := hex.EncodeToString(tx.Id)

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxId := hex.EncodeToString(in.TxId)
					spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
				}
			}

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[encodedTxId] != nil {
					for _, spentOut := range spentTXOs[encodedTxId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, tx)
				}
			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs, nil
}

func (blockchain *Blockchain) FindUTXO(address string) ([]*TXOutput, error) {
	var UTXOs []*TXOutput
	unspent, err := blockchain.FindUnspentTransactions(address)
	if err != nil {
		return nil, err
	}

	for _, tx := range unspent {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs, nil
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: bc.Tip,
		DB:          bc.DB,
	}
}

func (blockchain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int, error) {
	unspentOutputs := make(map[string][]int)
	unspentTXs, err := blockchain.FindUnspentTransactions(address)
	if err != nil {
		return 0, nil, err
	}

	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		encodedTXId := hex.EncodeToString(tx.Id)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[encodedTXId] = append(unspentOutputs[encodedTXId], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs, nil
}
