package passphrase

import "github.com/jedib0t/go-passwords/passphrase/dictionaries"

// Rule controls how the Generator/Sequencer generates passwords.
type Rule func(a any)

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
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.capitalize = enabled
		}
	}
}

func WithDictionary(words []string) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.dictionary = words
		}
	}
}

// WithNumber injects a random number after one of the words in the passphrase.
func WithNumber(enabled bool) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.withNumber = enabled
		}
	}
}

// WithNumWords sets the number of words in the passphrase.
func WithNumWords(n int) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.numWords = n
		}
	}
}

// WithSeparator sets up the delimiter to separate words.
func WithSeparator(s string) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.separator = s
		}
	}
}

func WithWordLength(min, max int) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.wordLenMin = min
			v.wordLenMax = max
		}
	}
}
