package rng

import (
	"testing"
)

func FuzzIntN(f *testing.F) {
	f.Add(-1)
	f.Add(0)
	f.Add(1)
	f.Add(2)
	f.Add(66)
	f.Add(129)
	f.Add(256)
	f.Add(257)
	f.Add(1 << 20)
	f.Add(1 << 33)
	f.Fuzz(func(t *testing.T, n int) {
		v, err := IntN(n)
		if n < 1 {
			if err != ErrInvalidN {
				t.Fatalf("IntN(%d) = (%d, %v), want ErrInvalidN", n, v, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("IntN(%d) failed: %v", n, err)
		}
		if v < 0 || v >= n {
			t.Fatalf("IntN(%d) = %d, out of range [0, %d)", n, v, n)
		}
	})
}

func FuzzFillIntNs(f *testing.F) {
	f.Add(1, 1)
	f.Add(66, 100)
	f.Add(256, 64)
	f.Add(1<<20, 16)
	f.Fuzz(func(t *testing.T, n int, count int) {
		if count < 0 || count > 4096 {
			t.Skip()
		}
		buf := make([]int, count)
		err := FillIntNs(buf, n)
		if n < 1 {
			if err != ErrInvalidN {
				t.Fatalf("FillIntNs(buf, %d) = %v, want ErrInvalidN", n, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("FillIntNs(buf, %d) failed: %v", n, err)
		}
		for i, v := range buf {
			if v < 0 || v >= n {
				t.Fatalf("FillIntNs(buf, %d): buf[%d] = %d, out of range [0, %d)", n, i, v, n)
			}
		}
	})
}
