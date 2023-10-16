package lib

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"testing"
)

func TestErc20ABI(t *testing.T) {
	_, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		t.Error("Error creating a parsed ABI interface")
	}
}
