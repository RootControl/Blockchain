package utils

import (
	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func HashPublicKey(pubKey []byte) []byte {
	publicHash := Base58Encode(pubKey)

	return publicHash
}
