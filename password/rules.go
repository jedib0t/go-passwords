package password

import "github.com/jedib0t/go-passwords/charset"

// Rule controls how the Generator generates passwords.
type Rule func(g *generator)

var (
	basicRules = []Rule{
		WithCharset(charset.AllChars),
		WithLength(12),
	}
)

// WithCharset sets the Charset the Generator can use.
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
// can appear in the password. Negative values are treated as 0.
func WithMinLowerCase(min int) Rule {
	if min < 0 {
		min = 0
	}
	return func(g *generator) {
		g.minLowerCase = min
	}
}

// WithMinUpperCase controls the minimum number of upper case characters that
// can appear in the password. Negative values are treated as 0.
func WithMinUpperCase(min int) Rule {
	if min < 0 {
		min = 0
	}
	return func(g *generator) {
		g.minUpperCase = min
	}
}

// WithNumSymbols controls the min/max number of symbols that can appear in the
// password. When this rule is not used, symbols in the charset are treated
// like any other character and may appear any number of times.
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
		g.symbolsConfigured = true
	}
}
