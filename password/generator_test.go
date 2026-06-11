package password

import (
	"testing"
	"unicode"

	"github.com/jedib0t/go-passwords/charset"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	g, err := NewGenerator(
		WithCharset(charset.AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(t, err)

	charsetMap := make(map[rune]bool)
	for _, r := range charset.AlphaNumeric.WithoutAmbiguity().WithoutDuplicates() {
		charsetMap[r] = true
	}

	for idx := 0; idx < 100; idx++ {
		password, err := g.Generate()
		assert.NoError(t, err)
		assert.NotEmpty(t, password)
		assert.Equal(t, 12, len(password), "password should be 12 characters long")
		// Verify all characters are from the charset
		for _, r := range password {
			assert.True(t, charsetMap[r], "password contains invalid character: %c", r)
		}
	}
}

func TestGenerator_Generate_WithAMixOfEverything(t *testing.T) {
	g, err := NewGenerator(
		WithCharset(charset.AllChars.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
		WithMinLowerCase(5),
		WithMinUpperCase(2),
		WithNumSymbols(1, 1),
	)
	assert.Nil(t, err)

	for idx := 0; idx < 100; idx++ {
		password, err := g.Generate()
		assert.NoError(t, err)
		assert.NotEmpty(t, password)
		assert.Equal(t, 12, len(password), "password should be 12 characters long")

		numLowerCase := len(filterRunes([]rune(password), unicode.IsLower))
		assert.True(t, numLowerCase >= 5, "password: %s, lower case count: %d", password, numLowerCase)
		numUpperCase := len(filterRunes([]rune(password), unicode.IsUpper))
		assert.True(t, numUpperCase >= 2, "password: %s, upper case count: %d", password, numUpperCase)
		numSymbols := len(filterRunes([]rune(password), charset.Symbols.Contains))
		assert.True(t, numSymbols == 1, "password: %s, symbol count: %d", password, numSymbols)
	}
}

func TestGenerator_Generate_WithSymbols(t *testing.T) {
	t.Run("min 0 max 3", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdef123456-+!@#$%").WithoutAmbiguity().WithoutDuplicates()),
			WithLength(12),
			WithNumSymbols(0, 3),
		)
		assert.Nil(t, err)

		for idx := 0; idx < 100; idx++ {
			password, err := g.Generate()
			assert.NoError(t, err)
			assert.NotEmpty(t, password)
			assert.Equal(t, 12, len(password), "password should be 12 characters long")

			numSymbols := getNumSymbols(password)
			assert.True(t, numSymbols >= 0 && numSymbols <= 3, "password: %s, symbol count: %d", password, numSymbols)
		}
	})

	t.Run("min X max 3", func(t *testing.T) {
		for _, x := range []int{0, 1, 2, 3} {
			g, err := NewGenerator(
				WithCharset(charset.Charset("abcdef123456-+!@#$%").WithoutAmbiguity().WithoutDuplicates()),
				WithLength(12),
				WithNumSymbols(x, 3),
			)
			assert.Nil(t, err)

			for idx := 0; idx < 100; idx++ {
				password, err := g.Generate()
				assert.NoError(t, err)
				assert.NotEmpty(t, password)

				numSymbols := getNumSymbols(password)
				assert.True(t, numSymbols >= x, password)
				assert.True(t, numSymbols <= 3, password)
			}
		}
	})

	t.Run("min X max X", func(t *testing.T) {
		for _, x := range []int{0, 4, 6, 8, 12} {
			g, err := NewGenerator(
				WithCharset(charset.Charset("abcdef123456-+!@#$%").WithoutAmbiguity().WithoutDuplicates()),
				WithLength(12),
				WithNumSymbols(x, x),
			)
			assert.Nil(t, err)

			for idx := 0; idx < 100; idx++ {
				password, err := g.Generate()
				assert.NoError(t, err)
				assert.NotEmpty(t, password)
				assert.Equal(t, x, getNumSymbols(password), password)
			}
		}
	})
}

func TestGenerator_numSymbolsToGenerate(t *testing.T) {
	minSymbols, maxSymbols := 0, 3

	g := &generator{
		minSymbols: minSymbols,
		maxSymbols: maxSymbols,
	}
	for idx := 0; idx < 10000; idx++ {
		numSymbols, err := g.numSymbolsToGenerate()
		assert.NoError(t, err)
		assert.True(t, numSymbols >= minSymbols, numSymbols)
		assert.True(t, numSymbols <= maxSymbols, numSymbols)
	}
}

func getNumSymbols(pw string) int {
	rsp := 0
	for _, r := range pw {
		if charset.Symbols.Contains(r) {
			rsp++
		}
	}
	return rsp
}

func TestNewGenerator_ErrorCases(t *testing.T) {
	t.Run("empty charset", func(t *testing.T) {
		g, err := NewGenerator(WithCharset(charset.Charset("")))
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmptyCharset, err)
	})

	t.Run("zero length", func(t *testing.T) {
		g, err := NewGenerator(WithLength(0))
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrZeroLenPassword, err)
	})

	t.Run("no lower case in charset", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("ABCDEF123456")),
			WithLength(12),
			WithMinLowerCase(5),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNoLowerCaseInCharset, err)
	})

	t.Run("no upper case in charset", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdef123456")),
			WithLength(12),
			WithMinUpperCase(5),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNoUpperCaseInCharset, err)
	})

	t.Run("no symbols in charset", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdefABCDEF123456")),
			WithLength(12),
			WithNumSymbols(1, 1),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNoSymbolsInCharset, err)
	})

	t.Run("min lower case too long", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdef")),
			WithLength(5),
			WithMinLowerCase(10),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMinLowerCaseTooLong, err)
	})

	t.Run("min upper case too long", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("ABCDEF")),
			WithLength(5),
			WithMinUpperCase(10),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMinUpperCaseTooLong, err)
	})

	t.Run("min symbols too long", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdef!@#")),
			WithLength(5),
			WithNumSymbols(10, 10),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMinSymbolsTooLong, err)
	})

	t.Run("requirements not met", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdefABCDEF!@#")),
			WithLength(5),
			WithMinLowerCase(3),
			WithMinUpperCase(3),
			WithNumSymbols(1, 1),
		)
		assert.Nil(t, g)
		assert.NotNil(t, err)
		assert.Equal(t, ErrRequirementsNotMet, err)
	})
}

func TestWithNumSymbols_EdgeCases(t *testing.T) {
	t.Run("negative min", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdef!@#")),
			WithLength(12),
			WithNumSymbols(-5, 3),
		)
		assert.NotNil(t, g)
		assert.Nil(t, err)
		// min should be sanitized to 0
		for i := 0; i < 100; i++ {
			pw, err := g.Generate()
			assert.NoError(t, err)
			numSymbols := getNumSymbols(pw)
			assert.True(t, numSymbols >= 0 && numSymbols <= 3, "password: %s, symbols: %d", pw, numSymbols)
		}
	})

	t.Run("negative max", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdef!@#")),
			WithLength(12),
			WithNumSymbols(0, -5),
		)
		assert.NotNil(t, g)
		assert.Nil(t, err)
		// max should be sanitized to 0, so no symbols
		for i := 0; i < 100; i++ {
			pw, err := g.Generate()
			assert.NoError(t, err)
			numSymbols := getNumSymbols(pw)
			assert.Equal(t, 0, numSymbols, "password: %s", pw)
		}
	})

	t.Run("min > max", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(charset.Charset("abcdef!@#")),
			WithLength(12),
			WithNumSymbols(5, 3),
		)
		assert.NotNil(t, g)
		assert.Nil(t, err)
		// min should be set to max (3)
		for i := 0; i < 100; i++ {
			pw, err := g.Generate()
			assert.NoError(t, err)
			numSymbols := getNumSymbols(pw)
			assert.True(t, numSymbols >= 3 && numSymbols <= 3, "password: %s, symbols: %d", pw, numSymbols)
		}
	})
}

func TestGenerator_Generate_MaxSymbolsExceedsLength(t *testing.T) {
	// regression test: maxSymbols larger than the space left in the password
	// used to pass sanitize() and panic with an index-out-of-range inside
	// fill() whenever the drawn symbol count exceeded the password length
	g, err := NewGenerator(
		WithCharset(charset.Charset("abcdef!@#$%^&*")),
		WithLength(4),
		WithNumSymbols(0, 8),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)

	for i := 0; i < 500; i++ {
		pw, err := g.Generate()
		assert.NoError(t, err)
		assert.Equal(t, 4, len(pw))
		assert.LessOrEqual(t, getNumSymbols(pw), 4)
	}
}

func TestNewGenerator_WithBasicRules(t *testing.T) {
	// Test that NewGenerator applies basicRules by default
	g, err := NewGenerator()
	assert.NotNil(t, g)
	assert.Nil(t, err)

	// Should generate passwords with default settings (AllChars, length 12)
	for i := 0; i < 10; i++ {
		pw, err := g.Generate()
		assert.NoError(t, err)
		assert.Equal(t, 12, len(pw))
		assert.NotEmpty(t, pw)
	}
}

func TestGenerator_Generate_SymbolsWithoutNumSymbolsRule(t *testing.T) {
	// regression test: without a WithNumSymbols rule, symbols in the charset
	// must still be eligible for every position; the old implementation
	// silently filled the whole password from non-symbol characters only
	g, err := NewGenerator(
		WithCharset(charset.AllChars),
		WithLength(12),
	)
	assert.Nil(t, err)

	sawSymbol := false
	for i := 0; i < 200 && !sawSymbol; i++ {
		pw, err := g.Generate()
		assert.NoError(t, err)
		sawSymbol = getNumSymbols(pw) > 0
	}
	assert.True(t, sawSymbol, "expected at least one symbol across 200 default-rule passwords")
}

func TestNewGenerator_AllSymbolCharsetNeedsFillers(t *testing.T) {
	// regression test: an all-symbol charset with a symbol quota below the
	// password length used to pass sanitize() and then fail every Generate()
	// call, since the remaining positions are filled from non-symbols
	g, err := NewGenerator(
		WithCharset(charset.Symbols),
		WithLength(4),
		WithNumSymbols(1, 1),
	)
	assert.Nil(t, g)
	assert.Equal(t, ErrNoNonSymbolsInCharset, err)

	// an all-symbol charset that is fully covered by the symbol quota is fine
	g, err = NewGenerator(
		WithCharset(charset.Symbols),
		WithLength(4),
		WithNumSymbols(4, 4),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)
	pw, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(pw))
}

func TestGenerator_Generate_SymbolQuotaWithoutSymbolsInCharset(t *testing.T) {
	// regression test: a 0-N symbol quota on a charset with no symbols used
	// to draw a non-zero symbol count and fail Generate with an RNG error
	g, err := NewGenerator(
		WithCharset(charset.Charset("abcdef")),
		WithLength(6),
		WithNumSymbols(0, 3),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)

	for i := 0; i < 100; i++ {
		pw, err := g.Generate()
		assert.NoError(t, err)
		assert.Equal(t, 6, len(pw))
		assert.Equal(t, 0, getNumSymbols(pw))
	}
}

func TestGenerator_Generate_NegativeMinimums(t *testing.T) {
	// regression test (found by FuzzNewGenerator): negative minimums used to
	// inflate the space available for symbols and overrun the working buffer
	g, err := NewGenerator(
		WithCharset(charset.Charset("ab!@#$%^&*")),
		WithLength(4),
		WithMinLowerCase(-83),
		WithMinUpperCase(-6),
		WithNumSymbols(4, 11),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)

	for i := 0; i < 200; i++ {
		pw, err := g.Generate()
		assert.NoError(t, err)
		assert.Equal(t, 4, len(pw))
	}
}
