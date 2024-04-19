package password

import (
	"fmt"
	"slices"
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
	var actualPasswords []string
	for idx := 0; idx < 100; idx++ {
		password := g.Generate()
		assert.NotEmpty(t, password)
		if idx < len(expectedPasswords) {
			actualPasswords = append(actualPasswords, password)
			assert.Equal(t, expectedPasswords[idx], password)
		}
	}
	if !slices.Equal(expectedPasswords, actualPasswords) {
		for _, pw := range actualPasswords {
			fmt.Printf("%#v,\n", pw)
		}
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
		"r*rnUqHeg5QP",
		"m1RNe4@eXuda",
		"tq@wKqhhTMAK",
		"r1PkMr&qta2t",
		"hsPv#wzGiChh",
		"uth#CZeag1o7",
		"FeFKFxxaf@cq",
		"jxVK^1sRis6z",
		"bVrPjBRC@bqy",
		"f$orrWDzVYjx",
	}
	var actualPasswords []string
	for idx := 0; idx < 100; idx++ {
		password := g.Generate()
		assert.NotEmpty(t, password)
		if idx < len(expectedPasswords) {
			actualPasswords = append(actualPasswords, password)
			assert.Equal(t, expectedPasswords[idx], password)
		}

		numLowerCase := len(filterRunes([]rune(password), unicode.IsLower))
		assert.True(t, numLowerCase >= 5, password)
		numUpperCase := len(filterRunes([]rune(password), unicode.IsUpper))
		assert.True(t, numUpperCase >= 2, password)
		numSymbols := len(filterRunes([]rune(password), Symbols.Contains))
		assert.True(t, numSymbols == 1, password)
	}
	if !slices.Equal(expectedPasswords, actualPasswords) {
		for _, pw := range actualPasswords {
			fmt.Printf("%#v,\n", pw)
		}
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
			"435c33b31--d",
			"a6d4--c12c#5",
			"44c@c5$@acac",
			"f-d32ab154-5",
			"%5de2%6bd52e",
			"@+26ab1$#65+",
			"b%ccd!d-!414",
			"$6e2b63f@4+@",
			"!e4b+dda6$@2",
			"2+4f-cc$ef64",
		}
		var actualPasswords []string
		for idx := 0; idx < 100; idx++ {
			password := g.Generate()
			assert.NotEmpty(t, password)
			if idx < len(expectedPasswords) {
				actualPasswords = append(actualPasswords, password)
				assert.Equal(t, expectedPasswords[idx], password)
			}

			numSymbols := getNumSymbols(password)
			assert.True(t, numSymbols >= 0 && numSymbols <= 3, password)
		}
		if !slices.Equal(expectedPasswords, actualPasswords) {
			for _, pw := range actualPasswords {
				fmt.Printf("%#v,\n", pw)
			}
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
