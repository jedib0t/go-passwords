package sequencer

import "github.com/jedib0t/go-passwords/charset"

// Rule controls how the Generator/Sequencer generates passwords.
type Rule func(s *sequencer)

var (
	basicRules = []Rule{
		WithCharset(charset.AlphaNumeric),
		WithLength(12),
	}
)

// WithCharset sets the Charset the Generator/Sequencer can use.
func WithCharset(c charset.Charset) Rule {
	return func(s *sequencer) {
		s.charset = []rune(c)
	}
}

// WithLength sets the length of the generated password.
func WithLength(l int) Rule {
	return func(s *sequencer) {
		s.numChars = l
	}
}
