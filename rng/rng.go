package rng

import (
	"encoding/binary"
)

// IntN returns a random integer in [0, n) using crypto/rand.
func IntN(n int) (int, error) {
	if n <= 1 {
		return 0, ErrInvalidN
	}

	// For small n, use modulo directly as bias is negligible.
	if n <= 256 {
		var b [1]byte
		if err := readBytesBuffered(b[:]); err != nil {
			return 0, err
		}
		return int(b[0]) % n, nil
	}

	// For larger n, use rejection sampling to avoid modulo bias.
	max := uint32((uint64(1) << 32) / uint64(n) * uint64(n))
	if max == 0 {
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

// IntNs returns a slice of random integers in [0, n) using crypto/rand.
// It uses batching to reduce mutex contention and stack-allocated buffers for
// small requests to minimize heap allocations.
func IntNs(n int, count int) ([]int, error) {
	if n <= 1 {
		return nil, ErrInvalidN
	}
	if count <= 0 {
		return nil, nil
	}

	res := make([]int, count)

	// For small n, use modulo directly as bias is negligible.
	// We use a small stack buffer to avoid heap allocation for the temporary byte slice.
	if n <= 256 {
		var stackBuf [64]byte
		var b []byte = stackBuf[:]
		if count > len(stackBuf) {
			b = make([]byte, count)
		} else {
			b = b[:count]
		}

		if err := readBytesBuffered(b); err != nil {
			return nil, err
		}
		for i := 0; i < count; i++ {
			res[i] = int(b[i]) % n
		}
		return res, nil
	}

	// For larger n, use rejection sampling to avoid modulo bias.
	max := uint32((uint64(1) << 32) / uint64(n) * uint64(n))
	if max == 0 {
		// Fallback for extremely large n where simple modulo is acceptable or max calculation overflows.
		var stackBuf [64]byte
		var b []byte = stackBuf[:]
		if count*4 > len(stackBuf) {
			b = make([]byte, count*4)
		} else {
			b = b[:count*4]
		}

		if err := readBytesBuffered(b); err != nil {
			return nil, err
		}
		for i := 0; i < count; i++ {
			res[i] = int(binary.BigEndian.Uint32(b[i*4:])) % n
		}
		return res, nil
	}

	// Rejection sampling loop to ensure zero bias.
	var b [4]byte
	for i := 0; i < count; i++ {
		for {
			if err := readBytesBuffered(b[:]); err != nil {
				return nil, err
			}
			val := binary.BigEndian.Uint32(b[:])
			if val < max {
				res[i] = int(val % uint32(n))
				break
			}
		}
	}
	return res, nil
}

// Shuffle shuffles the slice using Fisher-Yates algorithm with crypto/rand.
// For slices smaller than 256, it uses a batch of random bytes to avoid
// repeated RNG calls and mutex overhead.
func Shuffle[T any](slice []T) error {
	n := len(slice)
	if n <= 1 {
		return nil
	}

	// For small slices, batch the random bytes for all swaps (one byte per swap).
	if n <= 256 {
		var stackBuf [256]byte
		b := stackBuf[:n-1]
		if err := readBytesBuffered(b); err != nil {
			return err
		}
		for i := n - 1; i > 0; i-- {
			j := int(b[n-1-i]) % (i + 1)
			slice[i], slice[j] = slice[j], slice[i]
		}
		return nil
	}

	// For larger slices, fall back to individual IntN calls.
	for i := n - 1; i > 0; i-- {
		j, err := IntN(i + 1)
		if err != nil {
			return err
		}
		slice[i], slice[j] = slice[j], slice[i]
	}
	return nil
}
