package Services

import (
	"math/big"
	"Blockchain/src/Entities"
)

func CreateProofOfWork(block *Entities.Block) *Entities.ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Entities.TargetBits))
	
	pow := &Entities.ProofOfWork {
		Block: block,
		Target: target,
	}

	return pow
}