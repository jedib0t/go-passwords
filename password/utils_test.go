package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaximumPossibleWords(t *testing.T) {
	assert.Equal(t, "10", MaximumPossibleWords(Numbers, 1).String())
	assert.Equal(t, "10000", MaximumPossibleWords(Numbers, 4).String())
	assert.Equal(t, "100000000", MaximumPossibleWords(Numbers, 8).String())
	assert.Equal(t, "500246412961", MaximumPossibleWords(Symbols, 8).String())
	assert.Equal(t, "53459728531456", MaximumPossibleWords(Alphabets, 8).String())
	assert.Equal(t, "218340105584896", MaximumPossibleWords(AlphaNumeric, 8).String())
	assert.Equal(t, "4702525276151521", MaximumPossibleWords(AllChars, 8).String())
}
