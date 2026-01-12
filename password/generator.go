package password

import (
	"fmt"
	"sync"
	"unicode"
	"unicode/utf8"

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
	// GenerateTo generates a password and writes it to the provided buffer.
	// It returns the number of bytes written or an error.
	GenerateTo(buf []byte) (int, error)
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
			r := make([]rune, g.numChars)
			return &r
		},
	}
	for idx := 0; idx < storagePoolMinSize; idx++ {
		r := make([]rune, g.numChars)
		g.pool.Put(&r)
	}

	return g.sanitize()
}

// Generate returns a randomly generated password.
func (g *generator) Generate() (string, error) {
	buf := make([]byte, g.numChars*utf8.UTFMax)
	n, err := g.GenerateTo(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func (g *generator) GenerateTo(buf []byte) (int, error) {
	// use the pool to get a []rune for working on
	passwordPtr := g.pool.Get().(*[]rune)
	defer g.pool.Put(passwordPtr)
	password := *passwordPtr

	// fill it with minimum requirements first
	idx := 0
	if err := g.fill(password, g.charsetCaseLower, g.minLowerCase, &idx); err != nil {
		return 0, err
	}
	if err := g.fill(password, g.charsetCaseUpper, g.minUpperCase, &idx); err != nil {
		return 0, err
	}
	if numSymbols, err := g.numSymbolsToGenerate(); err != nil {
		return 0, err
	} else if err := g.fill(password, g.charsetSymbols, numSymbols, &idx); err != nil {
		return 0, err
	}
	if remainingChars := len(password) - idx; remainingChars > 0 {
		if err := g.fill(password, g.charsetNonSymbols, remainingChars, &idx); err != nil {
			return 0, err
		}
	}

	// shuffle it all
	if err := rng.Shuffle(password); err != nil {
		return 0, fmt.Errorf("failed to shuffle password: %w", err)
	}

	// write to the buffer
	return g.writeToBuf(password, buf)
}

func (g *generator) fill(password []rune, runes []rune, count int, idx *int) error {
	for ; *idx < len(password) && count > 0; count-- {
		n, err := rng.IntN(len(runes))
		if err != nil {
			return fmt.Errorf("failed to generate random number: %w", err)
		}
		password[*idx] = runes[n]
		(*idx)++
	}
	return nil
}

func (g *generator) writeToBuf(password []rune, buf []byte) (int, error) {
	offset := 0
	for _, r := range password {
		if offset+utf8.RuneLen(r) > len(buf) {
			return 0, ErrBufferTooSmall
		}
		offset += utf8.EncodeRune(buf[offset:], r)
	}
	return offset, nil
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
