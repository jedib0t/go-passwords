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
	charset    []rune
	charsetLen int
	numChars   int
	pool       *sync.Pool
	rng        *rand.Rand
}

// NewGenerator returns a password generator that implements the Generator
// interface.
func NewGenerator(charset Charset, numChars int) (Generator, error) {
	if len(charset) == 0 {
		return nil, ErrEmptyCharset
	}
	if numChars <= 0 {
		return nil, ErrZeroLenPassword
	}

	// create a storage pool with enough objects to support enough parallelism
	pool := &sync.Pool{
		New: func() any {
			return make([]rune, numChars)
		},
	}
	for idx := 0; idx < 25; idx++ {
		pool.Put(make([]rune, numChars))
	}

	return &generator{
		charset:    []rune(charset),
		charsetLen: len(charset),
		numChars:   numChars,
		pool:       pool,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
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
		// avoid repetition of previous character
		for idx > 0 && char == password[idx-1] {
			char = g.charset[g.rng.Intn(g.charsetLen)]
		}
		// set
		password[idx] = char
	}
	return string(password)
}

// SetSeed overrides the seed value for the RNG.
func (g *generator) SetSeed(seed int64) {
	g.rng = rand.New(rand.NewSource(seed))
}
