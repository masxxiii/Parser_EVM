package utils

import (
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	if strings.Compare(DB, "DB_LINK") != 0 {
		t.Errorf("got %q, wanted %q", DB, "DB_LINK")
	}

	if strings.Compare(RPC, "RPC_LINK") != 0 {
		t.Errorf("got %q, wanted %q", RPC, "RPC_LINK")
	}
}
