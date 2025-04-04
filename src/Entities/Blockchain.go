package Entities

type Blockchain struct {
	Blocks []*Block
}

// TODO: alter function to only add a type Block without knowing the previous block Hash
func (blockchain *Blockchain) AddBlock(newBlock *Block) {
	previousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	newBlock.PreviousBlockHash = previousBlock.Hash

	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}