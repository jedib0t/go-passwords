package passphrase

import (
	"fmt"
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
	Generate() string
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
func (g *generator) Generate() string {

	// select words
	var words []string
	var wordsMap = make(map[string]bool)
	for idx := 0; idx < g.numWords; idx++ {
		var word string
		for word == "" || wordsMap[word] {
			word = g.dictionary[rng.IntN(g.dictionaryLen)]
		}
		words = append(words, word)
		wordsMap[word] = true
	}

	// inject a random number after one of the words
	if g.withNumber {
		idx := rng.IntN(len(words))
		words[idx] += fmt.Sprint(rng.IntN(10))
	}

	return strings.Join(words, g.separator)
}

func (g *generator) sanitize() (Generator, error) {
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

	if g.numWords < NumWordsMin {
		return nil, ErrNumWordsTooSmall
	}
	if g.numWords > NumWordsMax {
		return nil, ErrNumWordsTooLarge
	}
	return g, nil
}
