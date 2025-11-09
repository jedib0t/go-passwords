package enumerator

import (
	"math/big"
)

func numValues(base int, length int) *big.Int {
	if length == 0 {
		return big.NewInt(0)
	}

	i, e := big.NewInt(int64(base)), big.NewInt(int64(length))
	return i.Exp(i, e, nil)
}
