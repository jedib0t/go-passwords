package rng

import (
	"encoding/binary"
)

// IntN returns a random integer in [0, n) using crypto/rand.
func IntN(n int) (int, error) {
	if n <= 1 {
		return 0, ErrInvalidN
	}

	// For small n, use single-byte rejection sampling to avoid modulo bias.
	if n <= 256 {
		limit := byteLimit(n)
		var b [1]byte
		for {
			if err := readBytesBuffered(b[:]); err != nil {
				return 0, err
			}
			if int(b[0]) < limit {
				return int(b[0]) % n, nil
			}
		}
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

// byteLimit returns the largest multiple of n that fits in [1, 256]. Bytes
// below this limit map uniformly onto [0, n) via modulo; bytes at or above it
// must be rejected and redrawn to avoid modulo bias.
func byteLimit(n int) int {
	return 256 - (256 % n)
}

// IntNs returns a slice of random integers in [0, n) using crypto/rand.
// It uses batching to reduce mutex contention and stack-allocated buffers for
// small requests to minimize heap allocations.
func IntNs(n int, count int) ([]int, error) {
	if count <= 0 {
		return nil, nil
	}
	res := make([]int, count)
	if err := FillIntNs(res, n); err != nil {
		return nil, err
	}
	return res, nil
}

// FillIntNs fills the provided slice with random integers in [0, n).
// It uses batching to reduce mutex contention and stack-allocated buffers for
// small requests to minimize additional heap allocations.
func FillIntNs(buf []int, n int) error {
	if n <= 1 {
		return ErrInvalidN
	}
	count := len(buf)
	if count <= 0 {
		return nil
	}

	// For small n, draw batches of bytes and filter them through rejection
	// sampling to avoid modulo bias. Batches are oversampled to cover the
	// expected rejection rate, and topped up until the buffer is full. A small
	// stack buffer avoids heap allocation for the temporary byte slice.
	if n <= 256 {
		limit := byteLimit(n)
		var stackBuf [128]byte
		filled := 0
		for filled < count {
			want := count - filled
			req := want + want/2 + 8 // oversample for rejected bytes
			if req > len(stackBuf) {
				req = len(stackBuf)
			}
			b := stackBuf[:req]
			if err := readBytesBuffered(b); err != nil {
				return err
			}
			for _, by := range b {
				if int(by) < limit {
					buf[filled] = int(by) % n
					filled++
					if filled == count {
						break
					}
				}
			}
		}
		return nil
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
			return err
		}
		for i := 0; i < count; i++ {
			buf[i] = int(binary.BigEndian.Uint32(b[i*4:])) % n
		}
		return nil
	}

	// Rejection sampling loop to ensure zero bias.
	var b [4]byte
	for i := 0; i < count; i++ {
		for {
			if err := readBytesBuffered(b[:]); err != nil {
				return err
			}
			val := binary.BigEndian.Uint32(b[:])
			if val < max {
				buf[i] = int(val % uint32(n))
				break
			}
		}
	}
	return nil
}

// Shuffle shuffles the slice using Fisher-Yates algorithm with crypto/rand.
// For slices smaller than 256, it uses batches of random bytes to avoid
// repeated RNG calls and mutex overhead. Each swap index is drawn with
// rejection sampling so that every permutation is equally likely.
func Shuffle[T any](slice []T) error {
	n := len(slice)
	if n <= 1 {
		return nil
	}

	// For small slices, batch the random bytes for the swaps and reject the
	// bytes that would introduce modulo bias for the current swap range.
	if n <= 256 {
		var stackBuf [128]byte
		var avail []byte
		for i := n - 1; i > 0; i-- {
			limit := byteLimit(i + 1)
			for {
				if len(avail) == 0 {
					req := i + i/2 + 8 // remaining swaps, oversampled for rejections
					if req > len(stackBuf) {
						req = len(stackBuf)
					}
					if err := readBytesBuffered(stackBuf[:req]); err != nil {
						return err
					}
					avail = stackBuf[:req]
				}
				by := avail[0]
				avail = avail[1:]
				if int(by) < limit {
					j := int(by) % (i + 1)
					slice[i], slice[j] = slice[j], slice[i]
					break
				}
			}
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
