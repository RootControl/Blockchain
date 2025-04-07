package domain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/rootcontrol/blockchain/pkg/utils"
)

const targetBits = 24
const maxNonce = 1 << 32

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{
		Block:  block,
		Target: target,
	}

	return pow
}

func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte {
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			utils.IntToHex(pow.Block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.Block.Data)
	for nonce < maxNonce {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		}

		nonce++
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.Target) == -1
}