package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/app/services"
	"github.com/rootcontrol/blockchain/internal/domain"
	"github.com/rootcontrol/blockchain/internal/infra/storage"
)

type CLI struct {
	BlockchainService *services.BlockchainService
	Repository interfaces.BlockchainRepository
}

func main() {
	repository, _ := storage.NewBoltRepository()
	defer repository.Close()

	cli := NewCLI(repository)

	cli.Run()
}

func NewCLI(repository interfaces.BlockchainRepository) *CLI {
	return &CLI {
		//BlockchainService: services.NewBlockchainService(repository, ""),
		Repository: repository,
	}
}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	sendData := sendCmd.String("from", "", "Transaction from")
	sendTo := sendCmd.String("to", "", "Transaction to")
	sendAmount := sendCmd.Int("amount", 0, "Transaction amount")
	getBalanceData := getBalanceCmd.String("address", "", "Transaction address")
	createBlockchainData := createBlockchainCmd.String("address", "", "Transaction address")

	switch os.Args[1] {
	case "createblockchain":
		createBlockchainCmd.Parse(os.Args[2:])
	case "send":
		sendCmd.Parse(os.Args[2:])
	case "getbalance":
		getBalanceCmd.Parse(os.Args[2:])
	case "createwallet":
		createWalletCmd.Parse(os.Args[2:])
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	if sendCmd.Parsed() {
		if *sendData == "" || *sendTo == "" || *sendAmount == 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.Send(*sendData, *sendTo, *sendAmount)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.GetBalance(*getBalanceData)
	}

	if createWalletCmd.Parsed() {
		wallet := domain.NewWallet()
		cli.Repository.SaveWallet(wallet)
		fmt.Printf("%s\n", wallet.GetAddress())
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainData == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}

		cli.BlockchainService = services.NewBlockchainService(cli.Repository, *createBlockchainData)
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
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - send AMOUNT of coins from FROM address to TO")
	fmt.Println("  getbalance -address ADDRESS - get balance for an address")
	fmt.Println("  createwallet - create a new wallet")
	fmt.Println("  createblockchain -address ADDRESS - create a blockchain and send genesis block reward to ADDRESS")
}

func (cli *CLI) Send(from, to string, amount int) {
	cli.BlockchainService = services.NewBlockchainService(cli.Repository, from)
	txService := services.NewTransactionService(cli.Repository, cli.BlockchainService.Blockchain.LastHash)
	transactions := txService.NewUnspentTxOutput(from, to, amount)
	
	err := cli.BlockchainService.MineBlock([]*domain.Transaction{transactions})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Success!")
}

func (cli *CLI) PrintChain() {
	iterator := services.NewIteratorService(
		cli.Repository,
		cli.BlockchainService.Blockchain.LastHash,
	)

	for {
		block := iterator.NextBlock()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) GetBalance(address string) {
	cli.BlockchainService = services.NewBlockchainService(cli.Repository, address)
	txService := services.NewTransactionService(cli.Repository, cli.BlockchainService.Blockchain.LastHash)
	balance := txService.GetBalance(address)

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
