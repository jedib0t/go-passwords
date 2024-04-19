package passphrase

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

const (
	MinWords = 256
)

type Generator interface {
	// Generate returns a randomly generated password.
	Generate() string
	// SetSeed overrides the seed value for the RNG.
	SetSeed(seed uint64)
}

type generator struct {
	capitalize bool
	dictionary []string
	separator  string
	numWords   int
	rng        *rand.Rand
	withNumber bool
	wordLenMin int
	wordLenMax int
}

// NewGenerator returns a password generator that implements the Generator
// interface.
func NewGenerator(rules ...Rule) (Generator, error) {
	g := &generator{}
	g.SetSeed(uint64(time.Now().UnixNano()))
	for _, opt := range append(defaultRules, rules...) {
		opt(g)
	}
	return g.sanitize()
}

// Generate returns a randomly generated password.
func (g *generator) Generate() string {
	var words []string

	// generate words
	for idx := 0; idx < g.numWords; idx++ {
		var word string
		for word == "" || slices.Contains(words, word) {
			word = g.dictionary[g.rng.IntN(len(g.dictionary))]
		}
		words = append(words, word)
	}
	// capitalize all words
	if g.capitalize {
		for idx := range words {
			words[idx] = strings.Title(words[idx])
		}
	}
	// inject a random number after one of the words
	if g.withNumber {
		idx := g.rng.IntN(len(words))
		words[idx] += fmt.Sprint(g.rng.IntN(10))
	}

	return strings.Join(words, g.separator)
}

// SetSeed overrides the seed value for the RNG.
func (g *generator) SetSeed(seed uint64) {
	g.rng = rand.New(rand.NewPCG(seed, seed+100))
}

func (g *generator) sanitize() (Generator, error) {
	if g.wordLenMin < 1 || g.wordLenMin > g.wordLenMax {
		return nil, ErrWordLengthInvalid
	}
	// filter the dictionary and remove too-short or too-long words
	slices.DeleteFunc(g.dictionary, func(word string) bool {
		return len(word) < g.wordLenMin || len(word) > g.wordLenMax
	})
	if len(g.dictionary) < MinWords {
		return nil, ErrDictionaryTooSmall
	}
	if g.numWords <= 0 {
		return nil, ErrNumWordsInvalid
	}
	return g, nil
}
