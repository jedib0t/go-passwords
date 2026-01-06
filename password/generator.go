package password

import (
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
	Generate() string
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
func (g *generator) Generate() string {
	// use the pool to get a []rune for working on
	password := g.pool.Get().([]rune)
	defer g.pool.Put(password)

	// init the filler
	idx := 0
	fillPassword := func(runes []rune, count int) {
		for ; idx < len(password) && count > 0; count-- {
			password[idx] = runes[rng.IntN(len(runes))]
			idx++
		}
	}

	// fill it with minimum requirements first
	if g.minLowerCase > 0 {
		fillPassword(g.charsetCaseLower, g.minLowerCase)
	}
	if g.minUpperCase > 0 {
		fillPassword(g.charsetCaseUpper, g.minUpperCase)
	}
	if numSymbols := g.numSymbolsToGenerate(); numSymbols > 0 {
		fillPassword(g.charsetSymbols, numSymbols)
	}
	// fill the rest with non-symbols (as symbols has a max)
	if remainingChars := len(password) - idx; remainingChars > 0 {
		fillPassword(g.charsetNonSymbols, remainingChars)
	}

	// shuffle it all
	rng.Shuffle(password)

	return string(password)
}

func (g *generator) numSymbolsToGenerate() int {
	if g.minSymbols > 0 || g.maxSymbols > 0 {
		return rng.IntN(g.maxSymbols-g.minSymbols+1) + g.minSymbols
	}
	return 0
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
