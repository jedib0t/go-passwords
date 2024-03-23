package password

import (
	"math/rand"
	"strings"
)

// Charset contains the list of allowed characters to use for the password
// generation.
type Charset string

// Some well-defined Charsets.
const (
	AlphabetsLower Charset = "abcdefghijklmnopqrstuvwxyz"
	AlphabetsUpper Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers        Charset = "0123456789"
	Symbols        Charset = "!@#$%^&*()-_=+"

	AllChars     = AlphaNumeric + Symbols
	AlphaNumeric = Alphabets + Numbers
	Alphabets    = AlphabetsUpper + AlphabetsLower

	ambiguous Charset = "O0lI"
)

// Contains returns true if the Charset contains the given char/rune.
func (c Charset) Contains(r rune) bool {
	for _, r2 := range c {
		if r == r2 {
			return true
		}
	}
	return false
}

// Shuffle reorders the Charset using the given RNG.
func (c Charset) Shuffle(rng *rand.Rand) Charset {
	cRunes := []rune(c)
	rng.Shuffle(len(cRunes), func(i, j int) {
		cRunes[i], cRunes[j] = cRunes[j], cRunes[i]
	})
	return Charset(cRunes)
}

// WithoutAmbiguity removes Ambiguous looking characters.
func (c Charset) WithoutAmbiguity() Charset {
	sb := strings.Builder{}
	for _, r := range c {
		if strings.Contains(string(ambiguous), string(r)) {
			continue
		}
		sb.WriteRune(r)
	}
	return Charset(sb.String())
}

// WithoutDuplicates removes duplicate characters.
func (c Charset) WithoutDuplicates() Charset {
	sb := strings.Builder{}
	for _, r := range c {
		if strings.Contains(sb.String(), string(r)) {
			continue
		}
		sb.WriteRune(r)
	}
	return Charset(sb.String())
}
