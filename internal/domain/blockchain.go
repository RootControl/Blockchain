package domain

type Blockchain struct {
	LastHash []byte
}

func NewBlockchain() *Blockchain {
	return &Blockchain{}
}

func (bc *Blockchain) AddBlock(block *Block) {
	bc.LastHash = block.Hash
}
