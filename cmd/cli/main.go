package main

import (
	"fmt"
	"strconv"

	"github.com/rootcontrol/blockchain/internal/domain"
)

func main() {
	blockchain := domain.NewBlockchain()

	blockchain.AddBlock("Send 1 BTC to Ivan")
	blockchain.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range blockchain.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := domain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}