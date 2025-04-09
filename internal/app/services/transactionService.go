package services

import (
	"encoding/hex"

	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/domain"
	transactions "github.com/rootcontrol/blockchain/internal/domain/Transactions"
)

type TransactionService struct {
	Repository interfaces.BlockchainRepository
	CurrentHash []byte
}

func NewTransactionService(repo interfaces.BlockchainRepository, currentHash []byte) *TransactionService {
	return &TransactionService{
		Repository: repo,
		CurrentHash: currentHash,
	}
}

func (service *TransactionService) GetBalance(address string) int {
	balance := 0
	unspentTxOutputs := service.FindUnspentTxOutputs(address)

	for _, output := range unspentTxOutputs {
		balance += output.Value
	}

	return balance
}

func (service *TransactionService) FindUnspentTxOutputs(address string) []transactions.TxOutput {
	var unspentTxOutputs []transactions.TxOutput
	unspentTransactions := service.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, output := range tx.TxOutputs {
			if output.CanBeUnlockedWith(address) {
				unspentTxOutputs = append(unspentTxOutputs, *output)
			}
		}
	}

	return unspentTxOutputs
}

func (service *TransactionService) FindUnspentTransactions(address string) []domain.Transaction {
	var unspentTxs []domain.Transaction
	spentTxOutputs := make(map[string][]int)
	iterator := NewIteratorService(service.Repository, service.CurrentHash)

	for {
		block := iterator.NextBlock()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.Id)

		Outputs:
			for outIndex, output := range tx.TxOutputs {
				// Checks if the output was spented
				if spentTxOutputs[txId] != nil {
					for _, spentOutput := range spentTxOutputs[txId] {
						if spentOutput == outIndex {
							continue Outputs
						}
					}
				}

				if output.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, input := range tx.TxInputs {
					if input.CanUnlockOutputWith(address) {
						inputTxId := hex.EncodeToString(input.TxId)
						spentTxOutputs[inputTxId] = append(spentTxOutputs[inputTxId], input.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTxs
}