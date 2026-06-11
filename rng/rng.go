package rng

import (
	"encoding/binary"
	"math"
)

// IntN returns a random integer in [0, n) using crypto/rand. A value of 1 is
// valid and always yields 0; values below 1 return ErrInvalidN.
func IntN(n int) (int, error) {
	if n < 1 {
		return 0, ErrInvalidN
	}
	if n == 1 {
		return 0, nil
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

	// For n up to 2^32, use 4-byte rejection sampling to avoid modulo bias.
	if uint64(n) <= 1<<32 {
		limit := (uint64(1) << 32) / uint64(n) * uint64(n)
		var b [4]byte
		for {
			if err := readBytesBuffered(b[:]); err != nil {
				return 0, err
			}
			val := uint64(binary.BigEndian.Uint32(b[:]))
			if val < limit {
				return int(val % uint64(n)), nil
			}
		}
	}

	// For larger n, use 8-byte rejection sampling.
	return intN64(n)
}

// intN64 returns a random integer in [0, n) for n > 2^32. It reads 8 random
// bytes per attempt and rejects values at or above the largest multiple of n
// representable in 64 bits to avoid modulo bias.
func intN64(n int) (int, error) {
	un := uint64(n)
	rem := (math.MaxUint64%un + 1) % un // 2^64 mod n
	var b [8]byte
	for {
		if err := readBytesBuffered(b[:]); err != nil {
			return 0, err
		}
		val := binary.BigEndian.Uint64(b[:])
		if rem == 0 || val <= math.MaxUint64-rem {
			return int(val % un), nil
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
// small requests to minimize additional heap allocations. A value of 1 for n
// is valid and fills the slice with zeroes; values below 1 return ErrInvalidN.
func FillIntNs(buf []int, n int) error {
	if n < 1 {
		return ErrInvalidN
	}
	count := len(buf)
	if count <= 0 {
		return nil
	}
	if n == 1 {
		for i := range buf {
			buf[i] = 0
		}
		return nil
	}

	// For small n, use batched single-byte rejection sampling.
	if n <= 256 {
		return fillIntNsBytes(buf, n)
	}

	// For n up to 2^32, use 4-byte rejection sampling per value.
	if uint64(n) <= 1<<32 {
		limit := (uint64(1) << 32) / uint64(n) * uint64(n)
		var b [4]byte
		for i := 0; i < count; i++ {
			for {
				if err := readBytesBuffered(b[:]); err != nil {
					return err
				}
				val := uint64(binary.BigEndian.Uint32(b[:]))
				if val < limit {
					buf[i] = int(val % uint64(n))
					break
				}
			}
		}
		return nil
	}

	// For larger n, use 8-byte rejection sampling per value.
	for i := 0; i < count; i++ {
		v, err := intN64(n)
		if err != nil {
			return err
		}
		buf[i] = v
	}
	return nil
}

// fillIntNsBytes fills the provided slice with random integers in [0, n) for
// n <= 256. It draws batches of bytes and filters them through rejection
// sampling to avoid modulo bias. Batches are oversampled to cover the
// expected rejection rate, and topped up until the buffer is full. A small
// stack buffer avoids heap allocation for the temporary byte slice.
func fillIntNsBytes(buf []int, n int) error {
	count := len(buf)
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
