package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/app/services"
	"github.com/rootcontrol/blockchain/internal/infra/storage"
)

type CLI struct {
	BlockchainService *services.BlockchainService
	Repository *interfaces.BlockchainRepository
}

func main() {
	repository, _ := storage.NewBoltRepository()
	defer repository.Close()

	cli := NewCLI(repository)

	cli.Run()
}

func NewCLI(repository interfaces.BlockchainRepository) *CLI {
	return &CLI {
		BlockchainService: services.NewBlockchainService(repository),
		Repository: &repository,
	}
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
	fmt.Println("  printchain - print the blocks in the blockchain")
}

func (cli *CLI) AddBlock(data string) {
	err := cli.BlockchainService.AddBlock(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Success!")
}

func (cli *CLI) PrintChain() {
	iterator := services.NewIteratorService(
		*cli.Repository,
		cli.BlockchainService.Blockchain.LastHash,
	)

	for {
		block := iterator.NextBlock()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}