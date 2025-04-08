package transactions

const subsidy = 10

type Transaction struct {
	Id        []byte
	TxInputs  []*TxInput
	TxOutputs []*TxOutput
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == ""{
		data = "Reward to " + to
	}

	txInput := NewTxInput([]byte{}, -1, data)
	txOutput := NewTxOutput(subsidy, to)
	tx := NewTransaction(nil, []*TxInput{txInput}, []*TxOutput{txOutput})

	return tx
}

func NewTransaction(id []byte, txInputs []*TxInput, txOutputs []*TxOutput) *Transaction {
	return &Transaction{
		Id: id,
		TxInputs: txInputs,
		TxOutputs: txOutputs,
	}
}