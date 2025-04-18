package services

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"

	"github.com/rootcontrol/blockchain/internal/app/interfaces"
	"github.com/rootcontrol/blockchain/internal/domain"
	"github.com/rootcontrol/blockchain/pkg/utils"
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

func (service *TransactionService) FindUnspentTxOutputs(address string) []domain.TxOutput {
	var unspentTxOutputs []domain.TxOutput
	unspentTransactions := service.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, output := range tx.TxOutputs {
			outputPubHash := output.CreateLockingScript([]byte(address))
			if output.IsLockedWithKey(outputPubHash) {
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

				outputPubHash := output.CreateLockingScript([]byte(address))
				if output.IsLockedWithKey(outputPubHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, input := range tx.TxInputs {
					if input.UsesKey(utils.Base58Encode([]byte(address))) {
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

func (service *TransactionService) NewUnspentTxOutput(from, to string, amount int) *domain.Transaction {
	var inputs []*domain.TxInput
	var outputs []*domain.TxOutput

	accumulated, validOutputs := service.FindSpendableOutputs(from, amount)

	if accumulated < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txIndex, outs := range validOutputs {
		txId, _ := hex.DecodeString(txIndex)

		for _, out := range outs {
			input := domain.NewTxInput(txId, out, []byte(from))
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	output := domain.NewTxOutput(amount, nil)
	output.Lock([]byte(to))
	outputs = append(outputs, output)

	// Create new TxOutput for the change
	if accumulated > amount {
		out := domain.NewTxOutput(accumulated - amount, nil)
		out.Lock([]byte(from))
		outputs = append(outputs, out)
	}

	transaction := domain.NewTransaction(nil, inputs, outputs)
	transaction.SetId()
	
	wallet, err := service.Repository.GetWallet(from)
	if err != nil {
		log.Panic(err)
	}

	service.SignTransaction(transaction, wallet.PrivateKey)

	return transaction
}

func (service *TransactionService) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTxs := service.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txId := hex.EncodeToString(tx.Id)

		for outIndex, output := range tx.TxOutputs {
			outputPubHash := output.CreateLockingScript([]byte(address))
			if output.IsLockedWithKey(outputPubHash) && accumulated < amount {
				accumulated += output.Value
				unspentOutputs[txId] = append(unspentOutputs[txId], outIndex)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (service *TransactionService) FindTransaction(id []byte) (domain.Transaction, error) {
	iterator := NewIteratorService(service.Repository, service.CurrentHash)

	for {
		block := iterator.NextBlock()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.Id, id) {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return domain.Transaction{}, errors.New("transaction not found")
}

func (service *TransactionService) SignTransaction(tx *domain.Transaction, privateKey ecdsa.PrivateKey) {
	prevTxs := service.GetPreviousTransactions(tx)

	tx.Sign(privateKey, prevTxs)
}

func (service *TransactionService) VerifyTransaction(tx *domain.Transaction) bool {
	prevTxs := service.GetPreviousTransactions(tx)

	return tx.Verify(prevTxs)
}

func (service *TransactionService) GetPreviousTransactions(tx *domain.Transaction) map[string]domain.Transaction {
	prevTxs := make(map[string]domain.Transaction)

	for _, input := range tx.TxInputs {
		prevTx, err := service.FindTransaction(input.TxId)

		if err != nil {
			log.Panic(err)
		}

		prevTxs[hex.EncodeToString(prevTx.Id)] = prevTx
	}

	return prevTxs
}
