package main

import (
	"Blockchain/src/Services"
	"fmt"
	"strconv"
)

func main() {
	blockchain := Services.NewBlockchain()

	blockchain.AddBlock("Send 1 BTC to Ivan")
	blockchain.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range blockchain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PreviousBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		// Show if the block is valid
		pow := Services.CreateProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}