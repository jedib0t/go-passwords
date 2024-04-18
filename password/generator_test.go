package password

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerator_Generate(b *testing.B) {
	g, err := NewGenerator(
		WithCharset(AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(b, err)
	assert.NotEmpty(b, g.Generate())

	for idx := 0; idx < b.N; idx++ {
		_ = g.Generate()
	}
}

func TestGenerator_Generate(t *testing.T) {
	g, err := NewGenerator(
		WithCharset(AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(t, err)
	g.SetSeed(1)

	expectedPasswords := []string{
		"KkeonkPQHv4r",
		"sHL31fveTcKB",
		"MHnSCTtKBds2",
		"oEqJBeZ8Qmie",
		"G2CGWSDAQUuz",
		"RtwGPgyAq9tN",
		"3kPAu4cMxN8t",
		"FgWWYrjqnx19",
		"uCFmDFDAoLZY",
		"pMgNoVa9z5Vv",
	}
	sb := strings.Builder{}
	for idx := 0; idx < 100; idx++ {
		password := g.Generate()
		assert.NotEmpty(t, password)
		if idx < len(expectedPasswords) {
			assert.Equal(t, expectedPasswords[idx], password)
			if expectedPasswords[idx] != password {
				sb.WriteString(fmt.Sprintf("%#v,\n", password))
			}
		}
	}
	if sb.Len() > 0 {
		fmt.Println(sb.String())
	}
}

func TestGenerator_Generate_WithAMixOfEverything(t *testing.T) {
	g, err := NewGenerator(
		WithCharset(AllChars.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
		WithMinLowerCase(5),
		WithMinUpperCase(2),
		WithNumSymbols(1, 1),
	)
	assert.Nil(t, err)
	g.SetSeed(1)

	expectedPasswords := []string{
		"r{rnUqHeg5QP",
		"m1RNe4$eXuda",
		"tq%wKqhhTMAK",
		"r1PkMr@qta2t",
		"hsPv+wzGiChh",
		"uth<CZeag1o7",
		"FeFKFxxaf|cq",
		"jxVK#1sRis6z",
		"bVrPjBRC<bqy",
		"f?orrWDzVYjx",
	}
	sb := strings.Builder{}
	for idx := 0; idx < 100; idx++ {
		password := g.Generate()
		assert.NotEmpty(t, password)
		if idx < len(expectedPasswords) {
			assert.Equal(t, expectedPasswords[idx], password)
			if expectedPasswords[idx] != password {
				sb.WriteString(fmt.Sprintf("%#v,\n", password))
			}
		}

		numLowerCase := len(filterRunes([]rune(password), unicode.IsLower))
		assert.True(t, numLowerCase >= 5, password)
		numUpperCase := len(filterRunes([]rune(password), unicode.IsUpper))
		assert.True(t, numUpperCase >= 2, password)
		numSymbols := len(filterRunes([]rune(password), Symbols.Contains))
		assert.True(t, numSymbols == 1, password)
	}
	if sb.Len() > 0 {
		fmt.Println(sb.String())
	}
}

func TestGenerator_Generate_WithSymbols(t *testing.T) {
	t.Run("min 0 max 3", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(Charset("abcdef123456-+!@#$%").WithoutAmbiguity().WithoutDuplicates()),
			WithLength(12),
			WithNumSymbols(0, 3),
		)
		assert.Nil(t, err)
		g.SetSeed(1)

		expectedPasswords := []string{
			"324c22b2f55c",
			"a5d355bf1c@4",
			"33c!b3#+acab",
			"e5c21aaf3353",
			"%3cd1$5bd31e",
			"+615abf$@536",
			"b%ccc-c5+3f3",
			"#4d1b52e!36!",
			"-e3a6cda4#!1",
			"162e6bb#ee53",
		}
		sb := strings.Builder{}
		for idx := 0; idx < 100; idx++ {
			password := g.Generate()
			assert.NotEmpty(t, password)
			if idx < len(expectedPasswords) {
				assert.Equal(t, expectedPasswords[idx], password)
				if expectedPasswords[idx] != password {
					sb.WriteString(fmt.Sprintf("%#v,\n", password))
				}
			}

			numSymbols := getNumSymbols(password)
			assert.True(t, numSymbols >= 0 && numSymbols <= 3, password)
		}
		if sb.Len() > 0 {
			fmt.Println(sb.String())
		}
	})

	t.Run("min X max 3", func(t *testing.T) {
		for _, x := range []int{0, 1, 2, 3} {
			g, err := NewGenerator(
				WithCharset(Charset("abcdef123456-+!@#$%").WithoutAmbiguity().WithoutDuplicates()),
				WithLength(12),
				WithNumSymbols(x, 3),
			)
			assert.Nil(t, err)
			g.SetSeed(1)

			for idx := 0; idx < 100; idx++ {
				password := g.Generate()
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
				WithCharset(Charset("abcdef123456-+!@#$%").WithoutAmbiguity().WithoutDuplicates()),
				WithLength(12),
				WithNumSymbols(x, x),
			)
			assert.Nil(t, err)
			g.SetSeed(1)

			for idx := 0; idx < 100; idx++ {
				password := g.Generate()
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
	g.SetSeed(1)
	for idx := 0; idx < 10000; idx++ {
		numSymbols := g.numSymbolsToGenerate()
		assert.True(t, numSymbols >= minSymbols, numSymbols)
		assert.True(t, numSymbols <= maxSymbols, numSymbols)
	}
}

func getNumSymbols(pw string) int {
	rsp := 0
	for _, r := range pw {
		if Symbols.Contains(r) {
			rsp++
		}
	}
	return rsp
}
