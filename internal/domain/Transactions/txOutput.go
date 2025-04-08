package transactions

type TxOutput struct {
	Value        int
	ScriptPubKey string
}

func NewTxOutput(value int, script string) *TxOutput {
	return &TxOutput{
		Value: value,
		ScriptPubKey: script,
	}
}