package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/joho/godotenv/autoload"
	"github.com/l-pay/evm_parser/core"
	"github.com/l-pay/evm_parser/db"
	"github.com/l-pay/evm_parser/lib"
	"github.com/l-pay/evm_parser/utils"
)

func startParser() {
	conn := db.ConnectDB()

	client, err := ethclient.Dial(os.Getenv(utils.RPC))
	if err != nil {
		log.Fatal(err)
	}

	erc20Abi, err := abi.JSON(strings.NewReader(lib.ERC20ABI))
	if err != nil {
		log.Fatal(err)
	}

	var startTime = time.Now()
	var sleepTime = time.Duration(10)

	for {
		utils.PrintServiceUptime(&startTime)

		parseErr := core.ProcessBlocks(client, conn, erc20Abi)
		if parseErr != nil {
			os.Exit(1)
		}

		utils.Sleep(&sleepTime)
	}
}

func main() {
	fmt.Println("\033[1;34m=========================================================================\033[0m")
	fmt.Println("\033[1;34m******************** [Ethereum] Block Parser Started ********************\033[0m")
	fmt.Println("\033[1;34m=========================================================================\033[0m")
	startParser()
}
