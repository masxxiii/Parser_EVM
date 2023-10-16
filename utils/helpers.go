package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

// PrintServiceUptime function prints the total uptime of the service since inception.
//   - Returns void
func PrintServiceUptime(startTime *time.Time) {
	fmt.Printf("\n\033[1;34m********************* Service Uptime [%012s] *********************\033[0m\n", time.Since(*startTime).Round(1*time.Second))
}

// PrintEndMessage function prints the end message after a batch of blocks is parsed
//   - Returns void
func PrintEndMessage() {
	fmt.Print("\n\u001B[1;36m*************************************************************************\u001B[0m")
	fmt.Println("\033[1;36m\n******************** WAITING FOR NEW BATCH OF BLOCKS ********************\033[0m")
	fmt.Print("\u001B[1;36m*************************************************************************\u001B[0m\n")
}

// Sleep function puts the system to sleep for a specified amount of time
//   - Returns void
func Sleep(seconds *time.Duration) {
	time.Sleep(*seconds * time.Second)
}

// TransferTopic 32 byte Keccak256 hash of transfer event
var TransferTopic = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
