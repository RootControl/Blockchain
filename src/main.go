package main

import (
	"Blockchain/src/Services"
	"fmt"
)

func main() {
	blockchain := Services.NewBlockchain()

	// TODO: alter way to add blocks
	Services.AddDataToBlockchain(blockchain, "Send 1 BTC to Ivan")
	Services.AddDataToBlockchain(blockchain, "Send 2 more BTC to Ivan")

	blocks := Services.GetAllBlocks(blockchain)

	for _, block := range blocks {
		fmt.Printf("Previous Hash: %x\n", block.PreviousBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}