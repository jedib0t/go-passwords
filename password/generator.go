package password

import (
	"fmt"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/jedib0t/go-passwords/charset"
	"github.com/jedib0t/go-passwords/rng"
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
	symbolsConfigured bool
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

	// create a storage pool for the working buffers; New covers misses, so
	// the pool needs no pre-filling
	g.pool = &sync.Pool{
		New: func() any {
			r := make([]rune, g.numChars)
			return &r
		},
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
	password := string(buf[:n])
	// wipe the intermediate buffer so the password lives only in the returned
	// string; callers needing full control of secret lifetime should use
	// GenerateTo with their own buffer instead
	clear(buf)
	return password, nil
}

func (g *generator) GenerateTo(buf []byte) (int, error) {
	// use the pool to get a []rune for working on; wipe it before returning
	// it so password material does not linger in pooled memory
	passwordPtr := g.pool.Get().(*[]rune)
	defer func() {
		clear(*passwordPtr)
		g.pool.Put(passwordPtr)
	}()
	password := (*passwordPtr)[:g.numChars]

	// fill it with minimum requirements first
	idx := 0
	if g.minLowerCase > 0 {
		if err := g.fill(password, g.charsetCaseLower, g.minLowerCase, &idx); err != nil {
			return 0, err
		}
	}
	if g.minUpperCase > 0 {
		if err := g.fill(password, g.charsetCaseUpper, g.minUpperCase, &idx); err != nil {
			return 0, err
		}
	}
	if numSymbols, err := g.numSymbolsToGenerate(); err != nil {
		return 0, err
	} else if numSymbols > 0 {
		if err := g.fill(password, g.charsetSymbols, numSymbols, &idx); err != nil {
			return 0, err
		}
	}
	if remainingChars := len(password) - idx; remainingChars > 0 {
		// when WithNumSymbols was configured, the symbol count generated above
		// is an exact quota, so the remaining characters must avoid symbols;
		// otherwise the full charset is fair game
		remainingCharset := g.charset
		if g.symbolsConfigured {
			remainingCharset = g.charsetNonSymbols
		}
		if err := g.fill(password, remainingCharset, remainingChars, &idx); err != nil {
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
	var stackBuf [64]int
	var indices []int = stackBuf[:]
	if count > len(stackBuf) {
		indices = make([]int, count)
	} else {
		indices = indices[:count]
	}

	if err := rng.FillIntNs(indices, len(runes)); err != nil {
		return fmt.Errorf("failed to generate random numbers: %w", err)
	}

	for _, n := range indices {
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
	if err := g.validateRules(); err != nil {
		return nil, err
	}
	// a symbol quota with no symbols in the charset can only ever yield zero
	// symbols; clamp it so numSymbolsToGenerate never asks for characters
	// that do not exist (minSymbols > 0 is already rejected above)
	if len(g.charsetSymbols) == 0 {
		g.maxSymbols = 0
	}
	// clamp maxSymbols to the space left over after the other minimums;
	// without this, numSymbolsToGenerate can exceed the password length and
	// overrun the working buffer
	if maxFit := g.numChars - g.minLowerCase - g.minUpperCase; g.maxSymbols > maxFit {
		g.maxSymbols = maxFit
	}
	return g, nil
}

func (g *generator) validateRules() error {
	if len(g.charset) == 0 {
		return ErrEmptyCharset
	}
	if g.numChars <= 0 {
		return ErrZeroLenPassword
	}
	if err := g.validateCaseRules(); err != nil {
		return err
	}
	if err := g.validateSymbolRules(); err != nil {
		return err
	}
	return nil
}

func (g *generator) validateCaseRules() error {
	if g.minLowerCase > 0 && len(g.charsetCaseLower) == 0 {
		return ErrNoLowerCaseInCharset
	}
	if g.minLowerCase > g.numChars {
		return ErrMinLowerCaseTooLong
	}
	if g.minUpperCase > 0 && len(g.charsetCaseUpper) == 0 {
		return ErrNoUpperCaseInCharset
	}
	if g.minUpperCase > g.numChars {
		return ErrMinUpperCaseTooLong
	}
	return nil
}

func (g *generator) validateSymbolRules() error {
	if g.minSymbols > 0 && len(g.charsetSymbols) == 0 {
		return ErrNoSymbolsInCharset
	}
	if g.minSymbols > g.numChars {
		return ErrMinSymbolsTooLong
	}
	if g.minLowerCase+g.minUpperCase+g.minSymbols > g.numChars {
		return ErrRequirementsNotMet
	}
	// when WithNumSymbols was configured, any characters beyond the drawn
	// symbol count are filled from the non-symbol charset; reject configs
	// that may need such fillers but have none available
	if g.symbolsConfigured && len(g.charsetNonSymbols) == 0 &&
		g.minLowerCase+g.minUpperCase+g.minSymbols < g.numChars {
		return ErrNoNonSymbolsInCharset
	}
	return nil
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
