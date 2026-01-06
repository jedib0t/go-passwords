package rng

import (
	"encoding/binary"
)

// IntN returns a random integer in [0, n) using crypto/rand.
func IntN(n int) (int, error) {
	if n <= 0 {
		return 0, ErrInvalidN
	}
	if n == 1 {
		return 0, ErrInvalidN
	}

	// For small n, use modulo directly (bias is negligible)
	// For larger n, use rejection sampling to avoid modulo bias
	if n <= 256 {
		var b [1]byte
		if err := readBytesBuffered(b[:]); err != nil {
			return 0, err
		}
		return int(b[0]) % n, nil
	}

	// Calculate the maximum value that is a multiple of n
	// to avoid modulo bias
	max := uint32((uint64(1) << 32) / uint64(n) * uint64(n))
	if max == 0 {
		// If max is 0, n is too large, fall back to simple modulo
		var b [4]byte
		if err := readBytesBuffered(b[:]); err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(b[:])) % n, nil
	}

	var b [4]byte
	for {
		if err := readBytesBuffered(b[:]); err != nil {
			return 0, err
		}
		val := binary.BigEndian.Uint32(b[:])
		if val < max {
			return int(val % uint32(n)), nil
		}
	}
}

// Shuffle shuffles the slice using Fisher-Yates algorithm with crypto/rand.
func Shuffle[T any](slice []T) error {
	for i := len(slice) - 1; i > 0; i-- {
		j, err := IntN(i + 1)
		if err != nil {
			return err
		}
		slice[i], slice[j] = slice[j], slice[i]
	}
	return nil
}
