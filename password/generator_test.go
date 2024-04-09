package password

import (
	"fmt"
	"strings"
	"testing"

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
		"7mfsqtsNrNcj",
		"LRagoTFi9iGz",
		"7RP42TsEV3se",
		"ZyL9CvRu5Ged",
		"y7y3wxVnRPMG",
		"qEJcmaT4yiL6",
		"XgVmEC15ZFH1",
		"XZzNALhgjLuV",
		"WoD3jbLU92tn",
		"XZxQd37ftWbo",
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
			"f4e23ab@6ecf",
			"2+2$e32fd3ce",
			"5#bd1bcfd12b",
			"3d$12$4%1c54",
			"ecfb5aed%3da",
			"b1a6ecf-dc-d",
			"bfb#a643bece",
			"6efd$36-3f1@",
			"%5f321f564eb",
			"eb5d3ef5-ef-",
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

func getNumSymbols(pw string) int {
	rsp := 0
	for _, r := range pw {
		if Symbols.Contains(r) {
			rsp++
		}
	}
	return rsp
}
