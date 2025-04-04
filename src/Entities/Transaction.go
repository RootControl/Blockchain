package Entities

type Transaction struct {
	Id []byte
	Inputs []TransactionInput
	Outputs []TransactionOutput
}

type TransactionInput struct {
	TransactionId []byte
	OutputIndex int
	Signature string
}

type TransactionOutput struct {
	Value int
	PubKey string
}

func (transaction *Transaction) SetId() {
	// transaction.Id = []byte{}
}