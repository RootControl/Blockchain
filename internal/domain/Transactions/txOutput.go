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

func (output *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return output.ScriptPubKey == unlockingData
}