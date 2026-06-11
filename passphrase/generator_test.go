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
		passphrase, err := g.Generate()
		assert.NoError(t, err)
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

func TestNewGenerator_DoesNotMutateCallerDictionary(t *testing.T) {
	// regression test: sanitize() filters, sorts, compacts and capitalizes
	// the dictionary; none of that may leak into the slice the caller passed
	words := []string{"zebra", "apple", "mango", "banana", "cherry", "papaya"}
	words = append(words, dictionaries.English()...)
	original := make([]string, len(words))
	copy(original, words)

	g, err := NewGenerator(
		WithDictionary(words),
		WithCapitalizedWords(true),
		WithNumWords(3),
		WithWordLength(4, 7),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)
	assert.Equal(t, original, words, "caller's dictionary slice must not be modified")
}

func TestGenerator_Generate_CapitalizationGrowsWordBytes(t *testing.T) {
	// regression test: 'ɐ' (U+0250, 2 bytes) capitalizes to 'Ɐ' (U+2C6F,
	// 3 bytes), so a word can exceed the word-length rule's byte limit after
	// capitalization; Generate used to size its buffer from the rule and
	// fail with ErrBufferTooSmall
	words := make([]string, 0, 300)
	for c1 := 'a'; c1 <= 'z'; c1++ {
		for c2 := 'a'; c2 <= 'z'; c2++ {
			words = append(words, "ɐ"+string(c1)+string(c2)+"def")
		}
	}

	g, err := NewGenerator(
		WithDictionary(words),
		WithCapitalizedWords(true),
		WithNumWords(3),
		WithNumber(true),
		WithSeparator("-"),
		WithWordLength(4, 7),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)

	for i := 0; i < 50; i++ {
		phrase, err := g.Generate()
		assert.NoError(t, err)
		assert.NotEmpty(t, phrase)
	}
}
