package password

import (
	"math/rand"
	"sync"
	"time"
)

var (
	storagePoolMinSize = 25
)

type Generator interface {
	// Generate returns a randomly generated password.
	Generate() string
	// SetSeed overrides the seed value for the RNG.
	SetSeed(seed int64)
}

type generator struct {
	charset           []rune
	charsetLen        int
	charsetSymbols    []rune
	charsetSymbolsLen int
	minSymbols        int
	maxSymbols        int
	numChars          int
	pool              *sync.Pool
	rng               *rand.Rand
}

// NewGenerator returns a password generator that implements the Generator
// interface.
func NewGenerator(rules ...Rule) (Generator, error) {
	g := &generator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for _, opt := range append(defaultRules, rules...) {
		opt(g)
	}

	// init the variables
	g.charsetLen = len(g.charset)
	g.charsetSymbols = []rune(Charset(g.charset).ExtractSymbols())
	g.charsetSymbolsLen = len(g.charsetSymbols)
	// create a storage pool with enough objects to support enough parallelism
	g.pool = &sync.Pool{
		New: func() any {
			return make([]rune, g.numChars)
		},
	}
	for idx := 0; idx < 25; idx++ {
		g.pool.Put(make([]rune, g.numChars))
	}

	// validate the inputs
	if g.charsetLen == 0 {
		return nil, ErrEmptyCharset
	}
	if g.numChars <= 0 {
		return nil, ErrZeroLenPassword
	}
	if g.minSymbols > g.numChars {
		return nil, ErrMinSymbolsTooLong
	}
	if g.minSymbols > 0 && g.charsetSymbolsLen == 0 {
		return nil, ErrNoSymbolsInCharset
	}
	return g, nil
}

// Generate returns a randomly generated password.
func (g *generator) Generate() string {
	// use the pool to get a []rune for working on
	password := g.pool.Get().([]rune)
	defer g.pool.Put(password)

	// overwrite the contents of the []rune and stringify it for response
	for idx := range password {
		// generate a random new character
		char := g.charset[g.rng.Intn(g.charsetLen)]

		// avoid repetition of previous character and ignore symbols
		for (idx > 0 && char == password[idx-1]) || Symbols.Contains(char) {
			char = g.charset[g.rng.Intn(g.charsetLen)]
		}

		// set
		password[idx] = char
	}

	// guarantee a minimum and maximum number of symbols
	if g.minSymbols > 0 || g.maxSymbols > 0 {
		numSymbolsToGenerate := g.minSymbols + g.rng.Intn(g.maxSymbols-g.minSymbols)
		for numSymbolsToGenerate > 0 {
			// generate a random new symbol
			char := g.charsetSymbols[g.rng.Intn(g.charsetSymbolsLen)]

			// find a random non-symbol location in the password
			location := g.rng.Intn(g.numChars)
			for Symbols.Contains(password[location]) {
				location = g.rng.Intn(g.numChars)
			}

			// set
			password[location] = char
			numSymbolsToGenerate--
		}
	}

	return string(password)
}

// SetSeed overrides the seed value for the RNG.
func (g *generator) SetSeed(seed int64) {
	g.rng = rand.New(rand.NewSource(seed))
}
