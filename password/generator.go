package password

import (
	"math/rand/v2"
	"sync"
	"time"
	"unicode"
)

var (
	storagePoolMinSize = 25
)

type Generator interface {
	// Generate returns a randomly generated password.
	Generate() string
	// SetSeed overrides the seed value for the RNG.
	SetSeed(seed uint64)
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
	rng               *rand.Rand
}

// NewGenerator returns a password generator that implements the Generator
// interface.
func NewGenerator(rules ...Rule) (Generator, error) {
	g := &generator{}
	g.SetSeed(uint64(time.Now().UnixNano()))
	for _, opt := range append(defaultRules, rules...) {
		opt(g)
	}

	// split the charsets
	g.charsetCaseLower = filterRunes(g.charset, unicode.IsLower)
	g.charsetCaseUpper = filterRunes(g.charset, unicode.IsUpper)
	g.charsetNonSymbols = filterRunes(g.charset, func(r rune) bool { return !Symbols.Contains(r) })
	g.charsetSymbols = filterRunes(g.charset, Symbols.Contains)

	// create a storage pool with enough objects to support enough parallelism
	g.pool = &sync.Pool{
		New: func() any {
			return make([]rune, g.numChars)
		},
	}
	for idx := 0; idx < 25; idx++ {
		g.pool.Put(make([]rune, g.numChars))
	}

	return g.sanitize()
}

// Generate returns a randomly generated password.
func (g *generator) Generate() string {
	// use the pool to get a []rune for working on
	password := g.pool.Get().([]rune)
	defer g.pool.Put(password)

	idx := 0
	fillPassword := func(count int, runes []rune) {
		for ; count > 0 && idx < len(password); count-- {
			password[idx] = runes[g.rng.IntN(len(runes))]
			idx++
		}
	}
	fillPassword(g.minLowerCase, g.charsetCaseLower)
	fillPassword(g.minUpperCase, g.charsetCaseUpper)
	fillPassword(g.numSymbolsToGenerate(), g.charsetSymbols)
	fillPassword(len(password)-idx, g.charsetNonSymbols)

	// shuffle it all
	g.rng.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}

// SetSeed overrides the seed value for the RNG.
func (g *generator) SetSeed(seed uint64) {
	g.rng = rand.New(rand.NewPCG(seed, seed+100))
}

func (g *generator) numSymbolsToGenerate() int {
	if g.minSymbols > 0 || g.maxSymbols > 0 {
		return g.rng.IntN(g.maxSymbols-g.minSymbols+1) + g.minSymbols
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
	if g.minSymbols > g.numChars {
		return nil, ErrMinSymbolsTooLong
	}
	if g.minLowerCase+g.minUpperCase+g.minSymbols > g.numChars {
		return nil, ErrRequirementsNotMet
	}
	if g.minSymbols > 0 && len(g.charsetSymbols) == 0 {
		return nil, ErrNoSymbolsInCharset
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
