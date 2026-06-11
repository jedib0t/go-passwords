package rng

import (
	"testing"
)

// chiSquared computes the chi-squared statistic of the observed counts
// against a uniform expectation over all buckets.
func chiSquared(counts []int, total int) float64 {
	expected := float64(total) / float64(len(counts))
	chi2 := 0.0
	for _, c := range counts {
		d := float64(c) - expected
		chi2 += d * d / expected
	}
	return chi2
}

// TestIntNUniformity guards against modulo bias in IntN. n=66 mirrors the
// default password charset size and is a worst case for single-byte modulo
// (256 % 66 = 58): the old biased implementation scores ~1400 on this
// statistic while a uniform source has df=65 (mean 65, stddev ~11.4). The
// bound of 250 is ~16 sigma and will not fail for a correct implementation.
func TestIntNUniformity(t *testing.T) {
	const n = 66
	const draws = 200_000
	counts := make([]int, n)
	for i := 0; i < draws; i++ {
		v, err := IntN(n)
		if err != nil {
			t.Fatal(err)
		}
		counts[v]++
	}
	if chi2 := chiSquared(counts, draws); chi2 > 250 {
		t.Errorf("IntN(%d) distribution is not uniform: chi-squared = %.1f (expected < 250)", n, chi2)
	}
}

// TestFillIntNsUniformity guards against modulo bias in the batched path of
// FillIntNs. Same statistic and bound as TestIntNUniformity.
func TestFillIntNsUniformity(t *testing.T) {
	const n = 66
	const draws = 200_000
	counts := make([]int, n)
	buf := make([]int, 1000)
	for i := 0; i < draws/len(buf); i++ {
		if err := FillIntNs(buf, n); err != nil {
			t.Fatal(err)
		}
		for _, v := range buf {
			counts[v]++
		}
	}
	if chi2 := chiSquared(counts, draws); chi2 > 250 {
		t.Errorf("FillIntNs(buf, %d) distribution is not uniform: chi-squared = %.1f (expected < 250)", n, chi2)
	}
}

// TestShuffleUniformity verifies that Shuffle produces every permutation of a
// small slice with equal probability. 4 elements have 24 permutations; with
// 120k shuffles each is expected 5000 times. df=23 (mean 23, stddev ~6.8);
// the bound of 120 is ~14 sigma.
func TestShuffleUniformity(t *testing.T) {
	const shuffles = 120_000
	counts := make([]int, 24)
	for i := 0; i < shuffles; i++ {
		s := []int{0, 1, 2, 3}
		if err := Shuffle(s); err != nil {
			t.Fatal(err)
		}
		counts[permIndex(s)]++
	}
	if chi2 := chiSquared(counts, shuffles); chi2 > 120 {
		t.Errorf("Shuffle permutation distribution is not uniform: chi-squared = %.1f (expected < 120)", chi2)
	}
}

// permIndex maps a permutation of {0,1,2,3} to a unique index in [0, 24)
// using the Lehmer code.
func permIndex(s []int) int {
	idx := 0
	factorials := []int{6, 2, 1}
	for i := 0; i < 3; i++ {
		smaller := 0
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[i] {
				smaller++
			}
		}
		idx += smaller * factorials[i]
	}
	return idx
}
