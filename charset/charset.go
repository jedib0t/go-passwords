package charset

import (
	"math/rand"
	"strings"

	crypto "github.com/jedib0t/go-passwords/rng"
)

// Charset contains the list of allowed characters to use for the password
// generation.
type Charset string

// Some well-defined Charsets.
const (
	AlphabetsLower Charset = "abcdefghijklmnopqrstuvwxyz"
	AlphabetsUpper Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers        Charset = "0123456789"
	Symbols        Charset = "!@#$%^&*"
	SymbolsFull    Charset = "~!@#$%^&*()-_=+[{]}|;:,<.>/?"

	AllChars     = AlphaNumeric + Symbols
	AlphaNumeric = Alphabets + Numbers
	Alphabets    = AlphabetsUpper + AlphabetsLower
)

var (
	ambiguousCharacters = map[rune]bool{
		'O': true,
		'0': true,
		'l': true,
		'I': true,
	}
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
//
// Deprecated: math/rand is not a cryptographically secure source of
// randomness; use Shuffled instead. Shuffling a charset does not by itself
// weaken generated passwords, but a seeded PRNG in a password library
// invites misuse.
func (c Charset) Shuffle(rng *rand.Rand) Charset {
	cRunes := []rune(c)
	rng.Shuffle(len(cRunes), func(i, j int) {
		cRunes[i], cRunes[j] = cRunes[j], cRunes[i]
	})
	return Charset(cRunes)
}

// Shuffled returns a copy of the Charset reordered using crypto/rand.
func (c Charset) Shuffled() (Charset, error) {
	cRunes := []rune(c)
	if err := crypto.Shuffle(cRunes); err != nil {
		return c, err
	}
	return Charset(cRunes), nil
}

// WithoutAmbiguity removes Ambiguous looking characters.
func (c Charset) WithoutAmbiguity() Charset {
	sb := strings.Builder{}
	for _, r := range c {
		if ambiguousCharacters[r] {
			continue
		}
		sb.WriteRune(r)
	}
	return Charset(sb.String())
}

// WithoutDuplicates removes duplicate characters.
func (c Charset) WithoutDuplicates() Charset {
	seen := make(map[rune]bool)
	sb := strings.Builder{}
	for _, r := range c {
		if seen[r] {
			continue
		}
		sb.WriteRune(r)
		seen[r] = true
	}
	return Charset(sb.String())
}
