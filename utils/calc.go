package utils

import (
	"math"
	"math/big"
	"strings"
)

// WeiToDecimal function to convert wei value to decimal representation based on decimals
//   - Returns decimal.Decimal
func WeiToDecimal(value string, decimals int) string {
	v, ok := new(big.Rat).SetString(value)
	if !ok {
		return ""
	}

	fromUnit := new(big.Int).SetInt64(int64(math.Pow10(0)))
	toUnit := new(big.Int).SetInt64(int64(math.Pow10(decimals)))

	return strings.TrimRight(v.Mul(v, new(big.Rat).SetFrac(fromUnit, toUnit)).FloatString(18), "0")
}
