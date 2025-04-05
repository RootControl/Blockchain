package Entities

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

const TargetBits = 24
const maxNonce = math.MaxInt64

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.Target) == -1
}

func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte {
			pow.Block.PreviousBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(TargetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte {},
	)

	return data
}

func IntToHex(num int64) []byte {
	return []byte(strconv.FormatInt(num, 16))
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			fmt.Printf("%x\n", hash)
			break
		}
		
		nonce++
	}

	return nonce, hash[:]
}