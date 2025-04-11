package domain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/rootcontrol/blockchain/pkg/utils"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet {
		PrivateKey: private,
		PublicKey: public,
	}

	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := utils.HashPublicKey(w.PublicKey)

	versionPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionPayload)

	fullPayload := append(versionPayload, checksum...)
	address := utils.Base58Encode(fullPayload)

	return address
}

func checksum(payload []byte) []byte {
	firstSha := sha256.Sum256(payload)
	secondSha := sha256.Sum256(firstSha[:])

	return secondSha[:addressChecksumLen]
}
