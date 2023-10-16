package core

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/l-pay/evm_parser/db"
	"github.com/l-pay/evm_parser/utils"
	"gorm.io/gorm"
)

// ProcessBlocks function to process the blocks that need to be parsed
//   - Returns error
func ProcessBlocks(eth *ethclient.Client, conn *gorm.DB, abi abi.ABI) error {
	latestBlock, BlockNumberErr := eth.BlockNumber(context.Background())
	if BlockNumberErr != nil {
		log.Fatal(BlockNumberErr)
		return BlockNumberErr
	}

	var localBlock = db.Block{ID: 1}
	err := conn.First(&localBlock).Error
	if err != nil {
		log.Fatal(err)
		return err
	}

	var diff = big.NewInt(0).Sub(big.NewInt(int64(latestBlock)), big.NewInt(int64(localBlock.Height)))

	if diff.Sign() == 1 {
		fmt.Printf("\033[1;34m********************* PARSING TOTAL OF [%03s] BLOCKS *********************\033[0m\n\n", diff.String())

		for b := int64(localBlock.Height) + 1; b <= int64(latestBlock); b++ {
			err := parseBlock(eth, conn, abi, big.NewInt(b))
			if err != nil {
				return err
			}
			fmt.Printf("\033[1;33m BLOCK [%010d] SUCCESSFULLY PARSED \033[0m\n", b)
		}

		utils.PrintEndMessage()
	}

	return nil
}

// parseBlock function to parse the given block number
//   - Returns error
func parseBlock(eth *ethclient.Client, conn *gorm.DB, abi abi.ABI, blockNumber *big.Int) error {
	query := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		Topics: [][]common.Hash{
			{utils.TransferTopic},
		},
	}

	block, err := eth.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
		return err
	}

	logs, err := eth.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	processNativeTransfers(block.Transactions(), conn)
	processTokenTransfers(logs, conn, abi)

	return nil
}

// processNativeTransfers function to process native transactions of the given block
//   - Returns void
func processNativeTransfers(transactions types.Transactions, conn *gorm.DB) {
	for _, trx := range transactions {
		if trx.To() == nil {
			continue
		}

		var to = strings.ToLower(trx.To().String())
		var wallet db.Wallet
		if err := conn.Where("address = ?", to).First(&wallet).Error; err != nil {
			continue
		}

		var hash = trx.Hash().Hex()
		var value = trx.Value().String()
		sender, err := types.Sender(types.LatestSignerForChainID(trx.ChainId()), trx)
		if err != nil {
			sender = common.HexToAddress("0x0000000000000000000000000000000000000000")
		}
		var from = strings.ToLower(sender.String())

		dbErr := addTransactionDB(conn, &hash, &from, &to, &value)
		if dbErr != nil {
			fmt.Printf("\u001B[1;31m [%s]\u001B[0m\n", hash)
			continue
		}

		fmt.Printf("\u001B[1;32m [%s]\u001B[0m\n", hash)
	}
}

// processTokenTransfers function to process ERC20 transfer events of the given block
//   - Returns void
func processTokenTransfers(logs []types.Log, conn *gorm.DB, abi abi.ABI) {
	var transferEvent struct {
		From  common.Address
		To    common.Address
		Value *big.Int
	}

	for _, l := range logs {
		err := abi.UnpackIntoInterface(&transferEvent, "Transfer", l.Data)
		if err != nil {
			continue
		}

		transferEvent.From = common.BytesToAddress(l.Topics[1].Bytes())
		transferEvent.To = common.BytesToAddress(l.Topics[2].Bytes())
		var to = strings.ToLower(transferEvent.To.String())
		var wallet db.Wallet
		if err := conn.Where("address = ?", to).First(&wallet).Error; err != nil {
			continue
		}

		var hash = l.TxHash.Hex()
		var value = transferEvent.Value.String()
		var from = strings.ToLower(transferEvent.From.String())

		dbErr := addTransactionDB(conn, &hash, &from, &to, &value)
		if dbErr != nil {
			fmt.Printf("\u001B[1;31m [%s]\u001B[0m\n", hash)
			continue
		}

		fmt.Printf("\u001B[1;32m [%s]\u001B[0m\n", hash)
	}
}

// addTransactionDB function to add parsed transaction to database
//   - Returns error
func addTransactionDB(conn *gorm.DB, hash *string, from *string, to *string, value *string) error {
	err := conn.Create(&db.Transaction{
		Hash:  *hash,
		From:  *from,
		To:    *to,
		Value: *value,
	}).Error

	if err != nil {
		return err
	}

	return nil
}
