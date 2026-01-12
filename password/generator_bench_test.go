package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerator_Generate(b *testing.B) {
	g, err := NewGenerator()
	assert.Nil(b, err)
	pw, err := g.Generate()
	assert.Nil(b, err)
	assert.NotEmpty(b, pw)

	for idx := 0; idx < b.N; idx++ {
		_, _ = g.Generate()
	}
}

func BenchmarkGenerator_GenerateTo(b *testing.B) {
	g, err := NewGenerator()
	assert.Nil(b, err)
	buf := make([]byte, 128)

	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		_, _ = g.GenerateTo(buf)
	}
}
