package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
)

const subsidy = 25

func NewCoinbaseTX(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := &TXInput{
		TxId:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}

	txout := &TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}

	tx := Transaction{
		Id:   nil,
		Vin:  []*TXInput{txin},
		Vout: []*TXOutput{txout},
	}

	if err := tx.SetId(); err != nil {
		return nil, err
	}

	return &tx, nil
}

type Transaction struct {
	Id   []byte
	Vin  []*TXInput
	Vout []*TXOutput
}

func (tx *Transaction) SetId() error {
	var (
		encoded bytes.Buffer
		hash    [32]byte
	)

	enc := gob.NewEncoder(&encoded)
	if err := enc.Encode(tx); err != nil {
		return err
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.Id = hash[:]

	return nil
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxId) == 0 && tx.Vin[0].Vout == -1
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (txOutput *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return txOutput.ScriptPubKey == unlockingData
}

type TXInput struct {
	TxId      []byte
	Vout      int
	ScriptSig string
}

func (txInput *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return txInput.ScriptSig == unlockingData
}

func NewUTXOTransaction(from, to string, amount int, blockchain *Blockchain) (*Transaction, error) {
	var (
		inputs  []*TXInput
		outputs []*TXOutput
	)

	acc, validOutputs, err := blockchain.FindSpendableOutputs(from, amount)
	if err != nil {
		return nil, err
	}

	if acc < amount {
		return nil, errors.New("Error: Not enough funds")
	}

	for txid, outs := range validOutputs {
		decodedTXId, err := hex.DecodeString(txid)
		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			inputs = append(inputs, &TXInput{
				TxId:      decodedTXId,
				Vout:      out,
				ScriptSig: from,
			})
		}
	}

	outputs = append(outputs, &TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, &TXOutput{
			Value:        acc - amount,
			ScriptPubKey: from,
		})
	}

	return &Transaction{
		Id:   nil,
		Vin:  inputs,
		Vout: outputs,
	}, nil
}
