package rng

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkIntN_Small(b *testing.B) {
	n := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := IntN(n)
		assert.Nil(b, err)
	}
}

func BenchmarkIntN_Medium(b *testing.B) {
	n := 10000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := IntN(n)
		assert.NoError(b, err)
	}
}

func BenchmarkIntN_Large(b *testing.B) {
	n := 1000000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := IntN(n)
		assert.NoError(b, err)
	}
}

func BenchmarkIntN_VeryLarge(b *testing.B) {
	n := 1 << 32 // 2^32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := IntN(n)
		assert.NoError(b, err)
	}
}

func BenchmarkShuffle_Small(b *testing.B) {
	slice := make([]rune, 10)
	for i := range slice {
		slice[i] = rune('a' + i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset slice for each iteration
		for j := range slice {
			slice[j] = rune('a' + j)
		}
		err := Shuffle(slice)
		assert.NoError(b, err)
	}
}

func BenchmarkShuffle_Medium(b *testing.B) {
	slice := make([]rune, 100)
	for i := range slice {
		slice[i] = rune(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset slice for each iteration
		for j := range slice {
			slice[j] = rune(j % 256)
		}
		err := Shuffle(slice)
		assert.NoError(b, err)
	}
}

func BenchmarkShuffle_Large(b *testing.B) {
	slice := make([]rune, 1000)
	for i := range slice {
		slice[i] = rune(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset slice for each iteration
		for j := range slice {
			slice[j] = rune(j % 256)
		}
		err := Shuffle(slice)
		assert.NoError(b, err)
	}
}
