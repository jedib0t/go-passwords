package password

import "github.com/jedib0t/go-passwords/charset"

// Rule controls how the Generator/Sequencer generates passwords.
type Rule func(g *generator)

var (
	basicRules = []Rule{
		WithCharset(charset.AllChars),
		WithLength(12),
	}
)

// WithCharset sets the Charset the Generator/Sequencer can use.
func WithCharset(c charset.Charset) Rule {
	return func(g *generator) {
		g.charset = []rune(c)
	}
}

// WithLength sets the length of the generated password.
func WithLength(l int) Rule {
	return func(g *generator) {
		g.numChars = l
	}
}

// WithMinLowerCase controls the minimum number of lower case characters that
// can appear in the password.
//
// Note: This works only on a Generator and is ineffective with a Sequencer.
func WithMinLowerCase(min int) Rule {
	return func(g *generator) {
		g.minLowerCase = min
	}
}

// WithMinUpperCase controls the minimum number of upper case characters that
// can appear in the password.
//
// Note: This works only on a Generator and is ineffective with a Sequencer.
func WithMinUpperCase(min int) Rule {
	return func(g *generator) {
		g.minUpperCase = min
	}
}

// WithNumSymbols controls the min/max number of symbols that can appear in the
// password.
//
// Note: This works only on a Generator and is ineffective with a Sequencer.
func WithNumSymbols(min, max int) Rule {
	// sanitize min and max
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if min > max {
		min = max
	}

	return func(g *generator) {
		g.minSymbols = min
		g.maxSymbols = max
	}
}
