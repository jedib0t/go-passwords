package password

import "math/big"

// MaximumPossibleWords returns the maximum number of unique passwords that can
// be generated with the given Charset and the number of characters allowed in
// the password.
func MaximumPossibleWords(charset Charset, numChars int) *big.Int {
	i, e := big.NewInt(int64(len(charset))), big.NewInt(int64(numChars))
	return i.Exp(i, e, nil)
}
