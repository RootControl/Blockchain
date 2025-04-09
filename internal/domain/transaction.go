package domain

import (
	"bytes"
	"encoding/gob"

	"crypto/sha256"

	"github.com/rootcontrol/blockchain/internal/domain/Transactions"
)

const subsidy = 10

type Transaction struct {
	Id        []byte
	TxInputs  []*transactions.TxInput
	TxOutputs []*transactions.TxOutput
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = "Reward to " + to
	}

	txInput := transactions.NewTxInput([]byte{}, -1, data)
	txOutput := transactions.NewTxOutput(subsidy, to)
	tx := NewTransaction(nil, []*transactions.TxInput{txInput}, []*transactions.TxOutput{txOutput})

	return tx
}

func NewTransaction(id []byte, txInputs []*transactions.TxInput, txOutputs []*transactions.TxOutput) *Transaction {
	return &Transaction{
		Id: id,
		TxInputs: txInputs,
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
