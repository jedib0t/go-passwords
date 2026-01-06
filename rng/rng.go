package rng

import (
	"crypto/rand"
	"encoding/binary"
)

// IntN returns a random integer in [0, n) using crypto/rand.
func IntN(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 0
	}

	// For small n, use modulo directly (bias is negligible)
	// For larger n, use rejection sampling to avoid modulo bias
	if n <= 256 {
		var b [1]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic(err)
		}
		return int(b[0]) % n
	}

	// Calculate the maximum value that is a multiple of n
	// to avoid modulo bias
	max := uint32((uint64(1) << 32) / uint64(n) * uint64(n))
	if max == 0 {
		// If max is 0, n is too large, fall back to simple modulo
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic(err)
		}
		return int(binary.BigEndian.Uint32(b[:])) % n
	}

	var b [4]byte
	for {
		if _, err := rand.Read(b[:]); err != nil {
			panic(err)
		}
		val := binary.BigEndian.Uint32(b[:])
		if val < max {
			return int(val % uint32(n))
		}
	}
}

// Shuffle shuffles the slice using Fisher-Yates algorithm with crypto/rand.
func Shuffle(slice []rune) {
	for i := len(slice) - 1; i > 0; i-- {
		j := IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
