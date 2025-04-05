package Services

import (
	"Blockchain/src/Entities"
	"fmt"
)

const subsidy = 10

func NewCoinbaseTransaction(to, data string) *Entities.Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	transactionInput := Entities.TransactionInput {
		TransactionId: []byte{},
		OutputIndex: -1,
		Signature: data,
	}

	transactionOutput := Entities.TransactionOutput {
		Value: subsidy,
		PubKey: to,
	}

	transaction := Entities.Transaction {
		Id: []byte{},
		Inputs: []Entities.TransactionInput{transactionInput},
		Outputs: []Entities.TransactionOutput{transactionOutput},
	}

	return &transaction
}