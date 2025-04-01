package Services

type Blockchain struct {
	Blocks []*Block
}

func (blockchain *Blockchain) AddBlock(data string) {
	previousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	newBlock := NewBlock(data, previousBlock.Hash)

	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}