package passphrase

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jedib0t/go-passwords/rng"
)

const (
	MinWordsInDictionary = 256
	NumWordsMin          = 2
	NumWordsMax          = 32
)

type Generator interface {
	// Generate returns a randomly generated password.
	Generate() (string, error)
}

type generator struct {
	capitalize    bool
	dictionary    []string
	dictionaryLen int
	separator     string
	numWords      int
	withNumber    bool
	wordLenMin    int
	wordLenMax    int
}

// NewGenerator returns a password generator that implements the Generator
// interface.
func NewGenerator(rules ...Rule) (Generator, error) {
	g := &generator{}
	for _, opt := range append(basicRules, rules...) {
		opt(g)
	}
	return g.sanitize()
}

// Generate returns a randomly generated password.
func (g *generator) Generate() (string, error) {
	// create and pre-allocate the builder
	var b strings.Builder
	// estimate the capacity needed
	b.Grow(g.wordLenMin*g.numWords + len(g.separator)*(g.numWords-1) + 1)

	// inject a random number after one of the words if asked for
	wordForDigitSuffixIdx, digit := -1, 0
	if g.withNumber {
		var err error

		wordForDigitSuffixIdx, err = rng.IntN(g.numWords)
		if err != nil {
			return "", err
		}
		digit, err = rng.IntN(10)
		if err != nil {
			return "", err
		}
	}

	// Select unique word indices using rejection sampling
	// For small numWords (typically 2-4), this is efficient and avoids large allocations
	wordIndicesMap := make(map[int]bool, g.numWords)

	// append words to the builder
	for idx := 0; idx < g.numWords; idx++ {
		// select a random word index from the dictionary (non-repeating)
		var wordIndex int
		var err error
		for wordIndex == 0 || wordIndicesMap[wordIndex] {
			wordIndex, err = rng.IntN(g.dictionaryLen)
			if err != nil {
				return "", err
			}
		}
		wordIndicesMap[wordIndex] = true
		// append the word to the builder
		b.WriteString(g.dictionary[wordIndex])

		// append the digit to the builder if asked for
		if wordForDigitSuffixIdx != -1 && idx == wordForDigitSuffixIdx {
			b.WriteString(string('0' + byte(digit)))
		}

		// append the separator if not the last word
		if idx < g.numWords-1 {
			b.WriteString(g.separator)
		}
	}

	return b.String(), nil
}

func (g *generator) sanitize() (Generator, error) {
	// check if the word length is valid
	if g.wordLenMin < 1 || g.wordLenMin > g.wordLenMax {
		return nil, ErrWordLengthInvalid
	}

	// remove words that are too-short & too-long
	g.dictionary = slices.DeleteFunc(g.dictionary, func(word string) bool {
		return len(word) < g.wordLenMin || len(word) > g.wordLenMax
	})
	slices.Sort(g.dictionary)
	g.dictionary = slices.Compact(g.dictionary)
	g.dictionaryLen = len(g.dictionary)

	// check if the dictionary is too small
	if g.dictionaryLen < g.numWords || g.dictionaryLen < MinWordsInDictionary {
		return nil, ErrDictionaryTooSmall
	}

	// capitalize all words in the dictionary ahead of time
	if g.capitalize {
		for idx := range g.dictionary {
			r, size := utf8.DecodeRuneInString(g.dictionary[idx])
			if r != utf8.RuneError {
				g.dictionary[idx] = string(unicode.ToUpper(r)) + g.dictionary[idx][size:]
			}
		}
	}

	// check if the number of words is too small or too large
	if g.numWords < NumWordsMin {
		return nil, ErrNumWordsTooSmall
	}
	if g.numWords > NumWordsMax {
		return nil, ErrNumWordsTooLarge
	}
	return g, nil
}
