package domain

import (
	"bytes"
	"encoding/gob"
	"crypto/sha256"
)

const subsidy = 10

type Transaction struct {
	Id        []byte
	TxInputs  []*TxInput
	TxOutputs []*TxOutput
}

func NewCoinbaseTx(to string) *Transaction {
	txInput := NewTxInput([]byte{}, -1, []byte{})
	txOutput := NewTxOutput(subsidy)
	txOutput.Lock([]byte(to))
	tx := NewTransaction(nil, []*TxInput{txInput}, []*TxOutput{txOutput})

	return tx
}

func NewTransaction(id []byte, txInputs []*TxInput, txOutputs []*TxOutput) *Transaction {
	return &Transaction{
		Id:        id,
		TxInputs:  txInputs,
		TxOutputs: txOutputs,
	}
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].TxId) == 0 && tx.TxInputs[0].Vout == -1
}

func (tx *Transaction) SetId() {
	txCopy := tx.Serialize()
	hash := sha256.Sum256(txCopy)
	tx.Id = hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	encoder.Encode(tx)

	return result.Bytes()
}
