package charset

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCharset(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	cs := AlphaNumeric.
		Shuffle(rng).
		WithoutAmbiguity().
		WithoutDuplicates()
	assert.Equal(t, "h72MfFSYvaW3GogBTczCysb4e5Ex9tXmpZHJUR8riAjLdNquPwkKQVDn6", string(cs))
}

func TestCharset_Contains(t *testing.T) {
	assert.True(t, Alphabets.Contains('a'))
	assert.True(t, AlphabetsLower.Contains('a'))
	assert.False(t, AlphabetsUpper.Contains('a'))
	assert.True(t, Symbols.Contains('#'))
}

func TestCharset_Shuffle(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	cs := Charset("abcde")
	cs = cs.Shuffle(rng)
	assert.Equal(t, "cdbae", string(cs))
	cs = cs.Shuffle(rng)
	assert.Equal(t, "baced", string(cs))
}

func TestCharset_WithoutAmbiguity(t *testing.T) {
	cs := Charset("abcde0oLlI1")
	cs = cs.WithoutAmbiguity()
	assert.Equal(t, "abcdeoL", string(cs))
}

func TestCharset_WithoutDuplicates(t *testing.T) {
	cs := Charset("abcde00oolI")
	cs = cs.WithoutDuplicates()
	assert.Equal(t, "abcde0olI", string(cs))
}

func TestCharset_Shuffled(t *testing.T) {
	cs := AlphaNumeric

	shuffled, err := cs.Shuffled()
	assert.NoError(t, err)
	assert.Len(t, shuffled, len(cs))
	for _, r := range cs {
		assert.True(t, shuffled.Contains(r), "shuffled charset should contain %c", r)
	}

	// with 62 characters, two crypto shuffles colliding with the original
	// order is impossible for all practical purposes
	shuffled2, err := cs.Shuffled()
	assert.NoError(t, err)
	assert.False(t, shuffled == cs && shuffled2 == cs, "shuffle should change the order")
}
