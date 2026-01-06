package password

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/jedib0t/go-passwords/charset"
	"github.com/jedib0t/go-passwords/rng"
)

var (
	// storagePoolMinSize is the minimum number of objects to keep in the pool
	// to support enough parallelism.
	storagePoolMinSize = 25
)

type Generator interface {
	// Generate returns a randomly generated password.
	Generate() (string, error)
}

type generator struct {
	charset           []rune
	charsetCaseLower  []rune
	charsetCaseUpper  []rune
	charsetNonSymbols []rune
	charsetSymbols    []rune
	minLowerCase      int
	minUpperCase      int
	minSymbols        int
	maxSymbols        int
	numChars          int
	pool              *sync.Pool
}

// NewGenerator returns a password generator that implements the Generator
// interface.
func NewGenerator(rules ...Rule) (Generator, error) {
	g := &generator{}
	for _, opt := range append(basicRules, rules...) {
		opt(g)
	}

	// split the charsets
	g.charsetCaseLower = filterRunes(g.charset, unicode.IsLower)
	g.charsetCaseUpper = filterRunes(g.charset, unicode.IsUpper)
	g.charsetNonSymbols = filterRunes(g.charset, func(r rune) bool { return !charset.Symbols.Contains(r) })
	g.charsetSymbols = filterRunes(g.charset, charset.Symbols.Contains)

	// create a storage pool with enough objects to support enough parallelism
	g.pool = &sync.Pool{
		New: func() any {
			return make([]rune, g.numChars)
		},
	}
	for idx := 0; idx < storagePoolMinSize; idx++ {
		g.pool.Put(make([]rune, g.numChars))
	}

	return g.sanitize()
}

// Generate returns a randomly generated password.
func (g *generator) Generate() (string, error) {
	// use the pool to get a []rune for working on
	password := g.pool.Get().([]rune)
	defer g.pool.Put(password)

	// init the filler
	idx := 0
	fillPassword := func(runes []rune, count int) error {
		for ; idx < len(password) && count > 0; count-- {
			n, err := rng.IntN(len(runes))
			if err != nil {
				return fmt.Errorf("failed to generate random number: %w", err)
			}
			password[idx] = runes[n]
			idx++
		}
		return nil
	}

	// fill it with minimum requirements first
	if g.minLowerCase > 0 {
		if err := fillPassword(g.charsetCaseLower, g.minLowerCase); err != nil {
			return "", err
		}
	}
	// fill the minimum upper case characters
	if g.minUpperCase > 0 {
		if err := fillPassword(g.charsetCaseUpper, g.minUpperCase); err != nil {
			return "", err
		}
	}
	// fill the minimum symbols
	if numSymbols, err := g.numSymbolsToGenerate(); err != nil {
		return "", fmt.Errorf("failed to generate number of symbols: %w", err)
	} else if numSymbols > 0 {
		if err := fillPassword(g.charsetSymbols, numSymbols); err != nil {
			return "", err
		}
	}
	// fill the rest with non-symbols (as symbols has a max)
	if remainingChars := len(password) - idx; remainingChars > 0 {
		if err := fillPassword(g.charsetNonSymbols, remainingChars); err != nil {
			return "", err
		}
	}

	// shuffle it all
	if err := rng.Shuffle(password); err != nil {
		return "", fmt.Errorf("failed to shuffle password: %w", err)
	}

	return string(password), nil
}

func (g *generator) numSymbolsToGenerate() (int, error) {
	if g.minSymbols > 0 || g.maxSymbols > 0 {
		// If min == max, no randomness needed
		if g.minSymbols == g.maxSymbols {
			return g.minSymbols, nil
		}
		n, err := rng.IntN(g.maxSymbols - g.minSymbols + 1)
		if err != nil {
			return 0, fmt.Errorf("failed to generate random number: %w", err)
		}
		return n + g.minSymbols, nil
	}
	return 0, nil
}

func (g *generator) sanitize() (Generator, error) {
	// validate the inputs
	if len(g.charset) == 0 {
		return nil, ErrEmptyCharset
	}
	if g.numChars <= 0 {
		return nil, ErrZeroLenPassword
	}
	if g.minLowerCase > 0 && len(g.charsetCaseLower) == 0 {
		return nil, ErrNoLowerCaseInCharset
	}
	if g.minLowerCase > g.numChars {
		return nil, ErrMinLowerCaseTooLong
	}
	if g.minUpperCase > 0 && len(g.charsetCaseUpper) == 0 {
		return nil, ErrNoUpperCaseInCharset
	}
	if g.minUpperCase > g.numChars {
		return nil, ErrMinUpperCaseTooLong
	}
	if g.minSymbols > 0 && len(g.charsetSymbols) == 0 {
		return nil, ErrNoSymbolsInCharset
	}
	if g.minSymbols > g.numChars {
		return nil, ErrMinSymbolsTooLong
	}
	if g.minLowerCase+g.minUpperCase+g.minSymbols > g.numChars {
		return nil, ErrRequirementsNotMet
	}
	return g, nil
}

func filterRunes(runes []rune, truth func(r rune) bool) []rune {
	var rsp []rune
	for _, r := range runes {
		if truth(r) {
			rsp = append(rsp, r)
		}
	}
	return rsp
}
