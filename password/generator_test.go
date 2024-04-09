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
			"f4e23abc6ecf",
			"242$e32fd3ce",
			"5ebd1bcfd12b",
			"213df12f4c1c",
			"54efb35ecfb5",
			"ed63@ad1eb1-",
			"6ecfbdcfd15b",
			"bfa643bece16",
			"fdf3$c3f1cba",
			"c345f321f56!",
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
	})

	t.Run("min 1 max 3", func(t *testing.T) {
		g, err := NewGenerator(
			WithCharset(Charset("abcdef123456-+!@#$%").WithoutAmbiguity().WithoutDuplicates()),
			WithLength(12),
			WithNumSymbols(1, 3),
		)
		assert.Nil(t, err)
		g.SetSeed(1)

		expectedPasswords := []string{
			"f4e23a%@6ecf",
			"424e-2fd#ce4",
			"ebd1bcf-12bc",
			"df1@f4c1c54@",
			"$cfb5aed%3da",
			"1a6@cfbdcfd1",
			"5bfbfa64-bec",
			"16efdf!6c3f1",
			"ac34#f321f5!",
			"@eb5d3ef5aef",
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
	})
}
