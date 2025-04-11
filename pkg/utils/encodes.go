package utils

import (
	"crypto/sha256"

	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func HashPublicKey(pubKey []byte) []byte {
	publicSha256 := sha256.Sum256(pubKey)

	hasher := sha256.New()
	hasher.Write(publicSha256[:])
	publicHash := hasher.Sum(nil)

	return publicHash
}