package main

import (
	"Blockchain/src/Services"
	"Blockchain/src/Infrastructures"
)

func main() {
	blockchain := Services.NewBlockchain()
	defer blockchain.Db.Close()

	cli := Infrastructures.CLI{Blockchain: blockchain}
	cli.Run()
}