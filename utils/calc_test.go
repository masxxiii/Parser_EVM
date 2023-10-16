package utils

import (
	"testing"
)

func TestWeiToDecimal(t *testing.T) {
	got := WeiToDecimal("1275328499080375895", 18)
	want := "1.275328499080375895"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
