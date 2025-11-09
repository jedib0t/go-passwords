package passphrase

import (
	"testing"

	"github.com/jedib0t/go-passwords/passphrase/dictionaries"
	"github.com/stretchr/testify/assert"
)

func TestNewGenerator_ErrorCases(t *testing.T) {
	t.Run("word length invalid - min < 1", func(t *testing.T) {
		g, err := NewGenerator(
			WithDictionary(dictionaries.English()),
			WithWordLength(0, 7),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrWordLengthInvalid, err)
	})

	t.Run("word length invalid - min > max", func(t *testing.T) {
		g, err := NewGenerator(
			WithDictionary(dictionaries.English()),
			WithWordLength(10, 5),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrWordLengthInvalid, err)
	})

	t.Run("dictionary too small after filtering", func(t *testing.T) {
		// Create a small dictionary that will be too small after filtering
		smallDict := []string{"a", "ab", "abc", "abcd", "abcde"}
		for i := 0; i < 250; i++ {
			smallDict = append(smallDict, "word")
		}

		g, err := NewGenerator(
			WithDictionary(smallDict),
			WithWordLength(10, 20), // This will filter out most words
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrDictionaryTooSmall, err)
	})

	t.Run("num words too small", func(t *testing.T) {
		g, err := NewGenerator(
			WithDictionary(dictionaries.English()),
			WithNumWords(1),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNumWordsTooSmall, err)
	})

	t.Run("num words too large", func(t *testing.T) {
		g, err := NewGenerator(
			WithDictionary(dictionaries.English()),
			WithNumWords(33),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNumWordsTooLarge, err)
	})
}

func TestGenerator_sanitize_EdgeCases(t *testing.T) {
	t.Run("capitalize with various word lengths", func(t *testing.T) {
		// Create a dictionary with enough words and various lengths
		dict := dictionaries.English()

		g, err := NewGenerator(
			WithDictionary(dict),
			WithCapitalizedWords(true),
			WithWordLength(4, 7),
		)
		assert.NotNil(t, g)
		assert.Nil(t, err)

		// Generate a few phrases to ensure it works
		for i := 0; i < 5; i++ {
			phrase := g.Generate()
			assert.NotEmpty(t, phrase)
		}
	})
}

func TestGenerator_Generate_EdgeCases(t *testing.T) {
	t.Run("without number", func(t *testing.T) {
		g, err := NewGenerator(
			WithDictionary(dictionaries.English()),
			WithNumWords(3),
			WithNumber(false),
			WithSeparator("-"),
		)
		assert.NotNil(t, g)
		assert.Nil(t, err)

		for i := 0; i < 10; i++ {
			phrase := g.Generate()
			assert.NotEmpty(t, phrase)
			// Should not contain digits
			hasDigit := false
			for _, r := range phrase {
				if r >= '0' && r <= '9' {
					hasDigit = true
					break
				}
			}
			assert.False(t, hasDigit, "phrase should not contain digits: %s", phrase)
		}
	})

	t.Run("with different separators", func(t *testing.T) {
		separators := []string{"-", "_", ".", " ", ""}
		for _, sep := range separators {
			g, err := NewGenerator(
				WithDictionary(dictionaries.English()),
				WithNumWords(2),
				WithSeparator(sep),
			)
			assert.NotNil(t, g)
			assert.Nil(t, err)

			phrase := g.Generate()
			assert.NotEmpty(t, phrase)
		}
	})
}
