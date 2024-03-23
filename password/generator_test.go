package password

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testGenCharset  = AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()
	testGenNumChars = 12
)

func BenchmarkGenerator_Generate(b *testing.B) {
	g, err := NewGenerator(testGenCharset, testGenNumChars)
	assert.Nil(b, err)
	assert.NotEmpty(b, g.Generate())

	for idx := 0; idx < b.N; idx++ {
		_ = g.Generate()
	}
}

func TestGenerator_Generate(t *testing.T) {
	g, err := NewGenerator(testGenCharset, testGenNumChars)
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
