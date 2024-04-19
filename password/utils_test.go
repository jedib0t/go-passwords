package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaximumPossibleWords(t *testing.T) {
	assert.Equal(t, "10", MaximumPossibleWords(Numbers, 1).String())
	assert.Equal(t, "10000", MaximumPossibleWords(Numbers, 4).String())
	assert.Equal(t, "100000000", MaximumPossibleWords(Numbers, 8).String())
	assert.Equal(t, "16777216", MaximumPossibleWords(Symbols, 8).String())
	assert.Equal(t, "377801998336", MaximumPossibleWords(SymbolsFull, 8).String())
	assert.Equal(t, "53459728531456", MaximumPossibleWords(Alphabets, 8).String())
	assert.Equal(t, "218340105584896", MaximumPossibleWords(AlphaNumeric, 8).String())
	assert.Equal(t, "576480100000000", MaximumPossibleWords(AllChars, 8).String())
	assert.Equal(t, "4304672100000000", MaximumPossibleWords(AlphaNumeric+SymbolsFull, 8).String())
}
