package sequencer

import (
	"testing"

	"github.com/jedib0t/go-passwords/charset"
	"github.com/stretchr/testify/assert"
)

func TestMaximumPossibleWords(t *testing.T) {
	assert.Equal(t, "10", MaximumPossibleWords(charset.Numbers, 1).String())
	assert.Equal(t, "10000", MaximumPossibleWords(charset.Numbers, 4).String())
	assert.Equal(t, "100000000", MaximumPossibleWords(charset.Numbers, 8).String())
	assert.Equal(t, "16777216", MaximumPossibleWords(charset.Symbols, 8).String())
	assert.Equal(t, "377801998336", MaximumPossibleWords(charset.SymbolsFull, 8).String())
	assert.Equal(t, "53459728531456", MaximumPossibleWords(charset.Alphabets, 8).String())
	assert.Equal(t, "218340105584896", MaximumPossibleWords(charset.AlphaNumeric, 8).String())
	assert.Equal(t, "576480100000000", MaximumPossibleWords(charset.AllChars, 8).String())
	assert.Equal(t, "4304672100000000", MaximumPossibleWords(charset.AlphaNumeric+charset.SymbolsFull, 8).String())
}
