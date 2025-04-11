package domain

import (
	"bytes"

	"github.com/rootcontrol/blockchain/pkg/utils"
)

type TxInput struct {
	TxId      []byte
	Vout      int
	Signature []byte
	PublicKey []byte
}

func NewTxInput(id []byte, vOut int, publicKey []byte) *TxInput {
	return &TxInput{
		TxId:      id,
		Vout:      vOut,
		PublicKey: publicKey,
	}
}

func (input *TxInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := utils.HashPublicKey(input.PublicKey)

	return bytes.Equal(lockingHash, publicKeyHash)
}