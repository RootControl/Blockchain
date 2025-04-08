package transactions

type TxInput struct {
	TxId      []byte
	Vout      int
	ScriptSig string
}

func NewTxInput(id []byte, vOut int, script string) *TxInput {
	return &TxInput{
		TxId: id,
		Vout: vOut,
		ScriptSig: script,
	}
}