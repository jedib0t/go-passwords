package passphrase

import "github.com/jedib0t/go-passwords/passphrase/dictionaries"

// Rule controls how the Generator/Sequencer generates passwords.
type Rule func(g *generator)

var (
	basicRules = []Rule{
		WithCapitalizedWords(true),
		WithDictionary(dictionaries.English()),
		WithNumWords(3),
		WithNumber(true),
		WithSeparator("-"),
		WithWordLength(4, 7),
	}
)

// WithCapitalizedWords ensures the words are Capitalized.
func WithCapitalizedWords(enabled bool) Rule {
	return func(g *generator) {
		g.capitalize = enabled
	}
}

// WithDictionary sets the dictionary of words to use for the passphrase.
func WithDictionary(words []string) Rule {
	return func(g *generator) {
		g.dictionary = words
	}
}

// WithNumber injects a random number after one of the words in the passphrase.
func WithNumber(enabled bool) Rule {
	return func(g *generator) {
		g.withNumber = enabled
	}
}

// WithNumWords sets the number of words in the passphrase.
func WithNumWords(n int) Rule {
	return func(g *generator) {
		g.numWords = n
	}
}

// WithSeparator sets up the delimiter to separate words.
func WithSeparator(s string) Rule {
	return func(g *generator) {
		g.separator = s
	}
}

// WithWordLength sets the minimum and maximum length of the words in the passphrase.
func WithWordLength(min, max int) Rule {
	return func(g *generator) {
		g.wordLenMin = min
		g.wordLenMax = max
	}
}
