package domain

type Blockchain struct {
	LastHash []byte
}

func NewBlockchain() *Blockchain {
	return &Blockchain{}
}

func (bc *Blockchain) SetLastHash(currentHash []byte) {
	bc.LastHash = currentHash
}
