package passphrase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerator_Generate(b *testing.B) {
	g, err := NewGenerator()
	assert.Nil(b, err)
	phrase, err := g.Generate()
	assert.Nil(b, err)
	assert.NotEmpty(b, phrase)

	for idx := 0; idx < b.N; idx++ {
		_, _ = g.Generate()
	}
}
