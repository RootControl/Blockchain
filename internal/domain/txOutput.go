package domain

import (
	"bytes"

	"github.com/rootcontrol/blockchain/pkg/utils"
)

type TxOutput struct {
	Value         int
	PublicKeyHash []byte
}

func NewTxOutput(value int, publicKeyHash []byte) *TxOutput {
	return &TxOutput{
		Value: value,
		PublicKeyHash: publicKeyHash,
	}
}

func (output *TxOutput) Lock(address []byte) {
	output.PublicKeyHash = output.CreateLockingScript(address)
}

func (output *TxOutput) CreateLockingScript(address []byte) []byte {
	publicKeyHash := utils.Base58Encode(address)
	publicKeyHash = publicKeyHash[1: len(publicKeyHash)-4]

	return publicKeyHash
}

func (output *TxOutput) IsLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Equal(output.PublicKeyHash, publicKeyHash)
}