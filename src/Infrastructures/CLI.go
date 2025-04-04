package Infrastructures

import (
	"Blockchain/src/Entities"
	"Blockchain/src/Services"
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	Blockchain *Entities.Blockchain
}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
		case "addblock":
	 		addBlockCmd.Parse(os.Args[2:])
		case "printchain":
			printChainCmd.Parse(os.Args[2:])
		default:
			cli.PrintUsage()
			os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}

		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) AddBlock(data string) {
	Services.AddDataToBlockchain(cli.Blockchain, data)
	fmt.Println("Success!")
}

func (cli *CLI) PrintChain() {
	blocks := cli.Blockchain.Iterator()

	for {
		block := Services.IterateNextBlock(blocks)

		if block == nil {
			break
		}

		fmt.Printf("Previous Hash: %x\n", block.PreviousBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}