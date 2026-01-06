package rng

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntN(t *testing.T) {
	t.Run("n <= 0", func(t *testing.T) {
		_, err := IntN(0)
		assert.Equal(t, ErrInvalidN, err)
		_, err = IntN(-1)
		assert.Equal(t, ErrInvalidN, err)
		_, err = IntN(-100)
		assert.Equal(t, ErrInvalidN, err)
	})

	t.Run("n == 1", func(t *testing.T) {
		result, err := IntN(1)
		assert.Equal(t, ErrInvalidN, err)
		assert.Equal(t, 0, result)
	})

	t.Run("small n (n <= 256)", func(t *testing.T) {
		// Test various small values
		for n := 2; n <= 256; n *= 2 {
			for i := 0; i < 100; i++ {
				result, err := IntN(n)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, result, 0, "result should be >= 0 for n=%d", n)
				assert.Less(t, result, n, "result should be < n for n=%d", n)
			}
		}
	})

	t.Run("medium n (256 < n < 2^32)", func(t *testing.T) {
		testCases := []int{257, 1000, 10000, 100000, 1000000, 10000000}
		for _, n := range testCases {
			for i := 0; i < 100; i++ {
				result, err := IntN(n)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, result, 0, "result should be >= 0 for n=%d", n)
				assert.Less(t, result, n, "result should be < n for n=%d", n)
			}
		}
	})

	t.Run("large n (n >= 2^32)", func(t *testing.T) {
		// Test with very large n that would cause max == 0
		largeN := 1 << 32 // 2^32
		for i := 0; i < 100; i++ {
			result, err := IntN(largeN)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, result, 0, "result should be >= 0 for large n")
			assert.Less(t, result, largeN, "result should be < n for large n")
		}
	})

	t.Run("range coverage", func(t *testing.T) {
		// Test that we get values across the entire range
		n := 10
		seen := make(map[int]bool)
		iterations := 1000
		for i := 0; i < iterations; i++ {
			result, err := IntN(n)
			assert.NoError(t, err)
			seen[result] = true
		}
		// With 1000 iterations and n=10, we should see all values
		assert.GreaterOrEqual(t, len(seen), n-1, "should see most values in range [0, %d)", n)
	})
}

func TestShuffle(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		slice := []rune{}
		original := make([]rune, len(slice))
		copy(original, slice)
		err := Shuffle(slice)
		assert.NoError(t, err)
		assert.Equal(t, original, slice, "empty slice should remain unchanged")
	})

	t.Run("single element", func(t *testing.T) {
		slice := []rune{'a'}
		original := make([]rune, len(slice))
		copy(original, slice)
		err := Shuffle(slice)
		assert.NoError(t, err)
		assert.Equal(t, original, slice, "single element slice should remain unchanged")
	})

	t.Run("two elements", func(t *testing.T) {
		slice := []rune{'a', 'b'}
		err := Shuffle(slice)
		assert.NoError(t, err)
		assert.Len(t, slice, 2, "slice should still have 2 elements")
		assert.Contains(t, slice, 'a', "slice should contain 'a'")
		assert.Contains(t, slice, 'b', "slice should contain 'b'")
	})

	t.Run("multiple elements", func(t *testing.T) {
		slice := []rune{'a', 'b', 'c', 'd', 'e'}
		original := make([]rune, len(slice))
		copy(original, slice)
		err := Shuffle(slice)
		assert.NoError(t, err)

		// Verify all elements are still present
		assert.Len(t, slice, len(original), "slice length should remain the same")
		for _, r := range original {
			assert.Contains(t, slice, r, "slice should contain all original elements")
		}

		// Count occurrences to ensure no duplicates
		counts := make(map[rune]int)
		for _, r := range slice {
			counts[r]++
		}
		for _, r := range original {
			assert.Equal(t, 1, counts[r], "each element should appear exactly once: %c", r)
		}
	})

	t.Run("shuffle actually changes order", func(t *testing.T) {
		// With a large enough slice and enough iterations,
		// we should see different orderings
		slice := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
		original := make([]rune, len(slice))
		copy(original, slice)

		allSame := true
		for i := 0; i < 100; i++ {
			copy(slice, original)
			err := Shuffle(slice)
			assert.NoError(t, err)
			if !equalSlices(slice, original) {
				allSame = false
				break
			}
		}
		assert.False(t, allSame, "shuffle should change the order of elements")
	})

	t.Run("large slice", func(t *testing.T) {
		slice := make([]rune, 1000)
		for i := range slice {
			slice[i] = rune(i % 256)
		}
		original := make([]rune, len(slice))
		copy(original, slice)
		err := Shuffle(slice)
		assert.NoError(t, err)

		// Verify all elements are still present (by counting)
		originalCounts := make(map[rune]int)
		for _, r := range original {
			originalCounts[r]++
		}
		shuffledCounts := make(map[rune]int)
		for _, r := range slice {
			shuffledCounts[r]++
		}
		assert.Equal(t, originalCounts, shuffledCounts, "element counts should match")
	})
}

func equalSlices(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestIntN_Concurrent(t *testing.T) {
	// Test that IntN works correctly under concurrent access
	// This verifies the buffering mechanism is thread-safe
	const goroutines = 100
	const iterations = 1000
	n := 100

	results := make(chan int, goroutines*iterations)
	errors := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					errors <- fmt.Errorf("panic: %v", r)
				}
			}()
			for j := 0; j < iterations; j++ {
				result, err := IntN(n)
				if err != nil {
					errors <- err
					return
				}
				if result < 0 || result >= n {
					errors <- fmt.Errorf("result %d out of range [0, %d)", result, n)
					return
				}
				results <- result
			}
		}()
	}

	// Collect all results
	seen := make(map[int]bool)
	for i := 0; i < goroutines*iterations; i++ {
		select {
		case result := <-results:
			seen[result] = true
		case err := <-errors:
			t.Fatal(err)
		}
	}

	// Verify we got a good distribution
	assert.GreaterOrEqual(t, len(seen), n/2, "should see at least half the range in concurrent access")
}
