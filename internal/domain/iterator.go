package domain

type BlockchainIterator struct {
	CurrentHash []byte
}

func NewBlockchainIterator(lastHash []byte) *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: lastHash,
	}
}
