package enumerator

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/jedib0t/go-passwords/charset"
	"github.com/stretchr/testify/assert"
)

func BenchmarkEnumerator_Decrement(b *testing.B) {
	o := New(charset.Numbers, 8, WithRolloverEnabled(true))

	for i := 0; i < b.N; i++ {
		_ = o.Decrement()
	}
}

func BenchmarkEnumerator_Decrement_Big(b *testing.B) {
	o := New(charset.AllChars, 256)
	o.Last()

	for i := 0; i < b.N; i++ {
		_ = o.Decrement()
	}
}

func BenchmarkEnumerator_DecrementN(b *testing.B) {
	o := New(charset.Numbers, 8, WithRolloverEnabled(true))

	n := big.NewInt(5)
	for i := 0; i < b.N; i++ {
		_ = o.DecrementN(n)
	}
}

func BenchmarkEnumerator_GoTo(b *testing.B) {
	o := New(charset.Numbers, 8, WithRolloverEnabled(true))
	maxValues := int64(math.Pow(10, 8))
	rng := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < b.N; i++ {
		// GoTo is 1-indexed, so generate values in [1, maxValues]
		n := big.NewInt(rng.Int63n(maxValues) + 1)
		err := o.GoTo(n)
		assert.Nil(b, err)
	}
}

func BenchmarkEnumerator_Increment(b *testing.B) {
	o := New(charset.Numbers, 8, WithRolloverEnabled(true))

	for i := 0; i < b.N; i++ {
		_ = o.Increment()
	}
}

func BenchmarkEnumerator_Increment_Big(b *testing.B) {
	o := New(charset.AllChars, 256)

	for i := 0; i < b.N; i++ {
		_ = o.Increment()
	}
}

func BenchmarkEnumerator_IncrementN(b *testing.B) {
	o := New(charset.Numbers, 8, WithRolloverEnabled(true))

	n := big.NewInt(5)
	for i := 0; i < b.N; i++ {
		_ = o.IncrementN(n)
	}
}

func BenchmarkEnumerator_String(b *testing.B) {
	o := New(charset.Numbers, 12)

	for i := 0; i < b.N; i++ {
		_ = o.String()
	}
}
