package password

// Rule controls how the Generator/Sequencer generates passwords.
type Rule func(any)

var (
	defaultRules = []Rule{
		WithCharset(AlphaNumeric),
		WithLength(8),
	}
)

// WithCharset sets the Charset the Generator/Sequencer can use.
func WithCharset(c Charset) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.charset = []rune(c)
		case *sequencer:
			v.charset = []rune(c)
		}
	}
}

// WithLength sets the length of the generated password.
func WithLength(l int) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.numChars = l
		case *sequencer:
			v.numChars = l
		}
	}
}

// WithMinLowerCase controls the minimum number of lower case characters that
// can appear in the password.
//
// Note: This works only on a Generator and is ineffective with a Sequencer.
func WithMinLowerCase(min int) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.minLowerCase = min
		}
	}
}

// WithMinUpperCase controls the minimum number of upper case characters that
// can appear in the password.
//
// Note: This works only on a Generator and is ineffective with a Sequencer.
func WithMinUpperCase(min int) Rule {
	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.minUpperCase = min
		}
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

	return func(a any) {
		switch v := a.(type) {
		case *generator:
			v.minSymbols = min
			v.maxSymbols = max
		}
	}
}
