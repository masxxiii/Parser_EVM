package utils

import (
	"testing"
)

func TestTransferTopic(t *testing.T) {
	got := TransferTopic
	want := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	if got.String() != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
