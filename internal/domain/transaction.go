package domain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"math/big"
)

const subsidy = 10

type Transaction struct {
	Id        []byte
	TxInputs  []*TxInput
	TxOutputs []*TxOutput
}

func NewCoinbaseTx(to string) *Transaction {
	txInput := NewTxInput([]byte{}, -1, []byte{})
	txOutput := NewTxOutput(subsidy, nil)
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

func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, previousTx map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for index, input := range txCopy.TxInputs {
		prevTx := previousTx[hex.EncodeToString(input.TxId)]
		txCopy.TxInputs[index].Signature = nil
		txCopy.TxInputs[index].PublicKey = prevTx.TxOutputs[input.Vout].PublicKeyHash

		txCopy.SetId()
		txCopy.TxInputs[index].PublicKey = nil

		r, s, _ := ecdsa.Sign(rand.Reader, &privateKey, txCopy.Id)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TxInputs[index].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput

	for _, input := range tx.TxInputs {
		inputs = append(inputs, NewTxInput(input.TxId, input.Vout, nil))
	}

	for _, output := range tx.TxOutputs {
		outputs = append(outputs, NewTxOutput(output.Value, output.PublicKeyHash))
	}

	txCopy := NewTransaction(tx.Id, inputs, outputs)

	return txCopy
}

func (tx *Transaction) Verify(previousTx map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for index, input := range tx.TxInputs {
		prevTx := previousTx[hex.EncodeToString(input.TxId)]
		txCopy.TxInputs[index].Signature = nil
		txCopy.TxInputs[index].PublicKey = prevTx.TxOutputs[input.Vout].PublicKeyHash
		txCopy.SetId()
		txCopy.TxInputs[index].PublicKey = nil
		
		r := big.Int{}
		s := big.Int{}
		sigLen := len(input.Signature)
		r.SetBytes(input.Signature[:(sigLen / 2)])
		s.SetBytes(input.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PublicKey)
		x.SetBytes(input.PublicKey[:(keyLen / 2)])
		y.SetBytes(input.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.Id, &r, &s) {
			return false
		}
	}	

	return true
}