package passphrase

import (
	"strings"
	"testing"

	"github.com/jedib0t/go-passwords/passphrase/dictionaries"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	g, err := NewGenerator(
		WithCapitalizedWords(true),
		WithDictionary(dictionaries.English()),
		WithNumWords(3),
		WithNumber(true),
		WithSeparator("-"),
		WithWordLength(4, 6),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)

	for idx := 0; idx < 1000; idx++ {
		passphrase := g.Generate()
		assert.NotEmpty(t, passphrase)

		// Verify structure: should have 3 words separated by "-"
		words := strings.Split(passphrase, "-")
		assert.Equal(t, 3, len(words), "passphrase should have 3 words: %s", passphrase)

		// Verify each word is capitalized and has a number in one of them
		hasNumber := false
		for _, word := range words {
			// Check that word starts with uppercase
			assert.True(t, len(word) > 0, "word should not be empty")
			assert.True(t, word[0] >= 'A' && word[0] <= 'Z', "word should start with uppercase: %s", word)
			// Check word length is between 4 and 6 (plus possibly a digit)
			assert.True(t, len(word) >= 4 && len(word) <= 7, "word length should be 4-7 (including possible digit): %s", word)

			// Check if word contains a digit
			for _, r := range word {
				if r >= '0' && r <= '9' {
					hasNumber = true
					break
				}
			}
		}
		assert.True(t, hasNumber, "passphrase should contain at least one number: %s", passphrase)
	}
}
