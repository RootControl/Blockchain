package domain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"

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

type serializableWallet struct {
	PrivateKeyD    *big.Int
	PrivateKeyX    *big.Int
	PrivateKeyY    *big.Int
	PrivateKeyCurve string
	PublicKey      []byte
}

func (w Wallet) Serialize() []byte {
	var result bytes.Buffer
	
	serializable := serializableWallet{
		PrivateKeyD:    w.PrivateKey.D,
		PrivateKeyX:    w.PrivateKey.PublicKey.X,
		PrivateKeyY:    w.PrivateKey.PublicKey.Y,
		PrivateKeyCurve: "P-256", // We're using P256 curve
		PublicKey:      w.PublicKey,
	}
	
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(serializable)
	
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeWallet(data []byte) *Wallet {
	var serializable serializableWallet

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&serializable)
	if err != nil {
		log.Panic(err)
	}

	privateKey := ecdsa.PrivateKey{
		D: serializable.PrivateKeyD,
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     serializable.PrivateKeyX,
			Y:     serializable.PrivateKeyY,
		},
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  serializable.PublicKey,
	}
}
