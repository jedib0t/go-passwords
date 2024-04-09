package password

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSequencer_GotoN(b *testing.B) {
	s, err := NewSequencer(
		WithCharset(AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(b, err)

	n := big.NewInt(math.MaxInt)
	pw, err := s.GotoN(n)
	assert.NotEmpty(b, pw)
	assert.Nil(b, err)
	assert.Equal(b, "AXZvFwUyHzQM", pw)

	for idx := 0; idx < b.N; idx++ {
		_, _ = s.GotoN(n)
	}
}

func BenchmarkSequencer_Next(b *testing.B) {
	s, err := NewSequencer(
		WithCharset(AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(b, err)
	s.First()

	assert.NotEmpty(b, s.Next())
	for idx := 0; idx < b.N; idx++ {
		s.Next()
	}
}

func BenchmarkSequencer_NextN(b *testing.B) {
	s, err := NewSequencer(
		WithCharset(AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(b, err)
	s.First()

	n := big.NewInt(100)
	assert.NotEmpty(b, s.NextN(n))
	for idx := 0; idx < b.N; idx++ {
		s.NextN(n)
	}
}

func BenchmarkSequencer_Prev(b *testing.B) {
	s, err := NewSequencer(
		WithCharset(AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(b, err)
	s.Last()

	assert.NotEmpty(b, s.Prev())
	for idx := 0; idx < b.N; idx++ {
		s.Prev()
	}
}

func BenchmarkSequencer_PrevN(b *testing.B) {
	s, err := NewSequencer(
		WithCharset(AlphaNumeric.WithoutAmbiguity().WithoutDuplicates()),
		WithLength(12),
	)
	assert.Nil(b, err)
	_, _ = s.GotoN(big.NewInt(math.MaxInt))

	n := big.NewInt(100)
	assert.NotEmpty(b, s.PrevN(n))
	for idx := 0; idx < b.N; idx++ {
		s.PrevN(n)
	}
}

func TestSequencer(t *testing.T) {
	s, err := NewSequencer(
		WithCharset(""),
		WithLength(3),
	)
	assert.Nil(t, s)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrEmptyCharset))
	s, err = NewSequencer(
		WithCharset("AB"),
		WithLength(0),
	)
	assert.Nil(t, s)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrZeroLenPassword))

	s, err = NewSequencer(
		WithCharset("AB"),
		WithLength(3),
	)
	assert.Nil(t, err)
	assert.Equal(t, "AAA", s.Get())
	assert.Equal(t, "0", s.GetN().String())
	assert.Equal(t, "AAA", s.First())
	assert.Equal(t, "0", s.GetN().String())
	assert.True(t, s.HasNext())
	assert.Equal(t, "BAA", s.NextN(big.NewInt(4)))
	assert.Equal(t, "4", s.GetN().String())
	assert.True(t, s.HasNext())
	assert.Equal(t, "BBA", s.NextN(big.NewInt(2)))
	assert.Equal(t, "6", s.GetN().String())
	assert.True(t, s.HasNext())
	assert.Equal(t, "BAA", s.PrevN(big.NewInt(2)))
	assert.Equal(t, "4", s.GetN().String())
	assert.True(t, s.HasNext())
	assert.Equal(t, "AAA", s.PrevN(big.NewInt(4)))
	assert.Equal(t, "0", s.GetN().String())
	assert.True(t, s.HasNext())
	assert.Equal(t, "BBB", s.Last())
	assert.Equal(t, "7", s.GetN().String())
	assert.False(t, s.HasNext())
	assert.Equal(t, "BBB", s.Get())

	// Next()
	expectedAnswers := []string{
		"AAA",
		"AAB",
		"ABA",
		"ABB",
		"BAA",
		"BAB",
		"BBA",
		"BBB",
	}
	s.Reset()
	assert.Equal(t, expectedAnswers[0], s.Get())
	for idx := 1; idx < len(expectedAnswers); idx++ {
		assert.Equal(t, expectedAnswers[idx], s.Next())
	}
	assert.Equal(t, "BBB", s.Next())
	assert.Equal(t, "BBB", s.NextN(big.NewInt(303)))

	// Prev()
	expectedAnswers = []string{
		"BBB",
		"BBA",
		"BAB",
		"BAA",
		"ABB",
		"ABA",
		"AAB",
		"AAA",
	}
	s.Last()
	assert.Equal(t, expectedAnswers[0], s.Get())
	for idx := 1; idx < len(expectedAnswers); idx++ {
		assert.Equal(t, expectedAnswers[idx], s.Prev())
	}
	assert.Equal(t, "AAA", s.Prev())
	assert.Equal(t, "AAA", s.PrevN(big.NewInt(303)))
}

func TestSequencer_GotoN(t *testing.T) {
	s, err := NewSequencer(
		WithCharset("AB"),
		WithLength(3),
	)
	assert.Nil(t, err)

	pw, err := s.GotoN(big.NewInt(-1))
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidN))
	pw, err = s.GotoN(big.NewInt(100))
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidN))

	s, err = NewSequencer(
		WithCharset("AB"),
		WithLength(4),
	)
	assert.Nil(t, err)
	expectedPasswords := []string{
		"AAAA",
		"AAAB",
		"AABA",
		"AABB",
		"ABAA",
		"ABAB",
		"ABBA",
		"ABBB",
		"BAAA",
		"BAAB",
		"BABA",
		"BABB",
		"BBAA",
		"BBAB",
		"BBBA",
		"BBBB",
	}
	for idx := 0; idx < len(expectedPasswords); idx++ {
		pw, err = s.GotoN(big.NewInt(int64(idx)))
		assert.Nil(t, err)
		assert.Equal(t, expectedPasswords[idx], pw, fmt.Sprintf("idx=%d", idx))
	}
}

func TestSequencer_Stream(t *testing.T) {
	s, err := NewSequencer(
		WithCharset("AB"),
		WithLength(3),
	)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ch := make(chan string, 1)
	// streams all possible passwords into the channel in an async routine
	go func() {
		t.Logf("< streaming passwords ...")
		err2 := s.Stream(ctx, ch)
		assert.Nil(t, err2)
		t.Logf("< streaming passwords ... done!")
	}()

	// listen on channel for passwords until channel is closed, or until timeout
	pw, ok := "", true
	var passwords []string
	t.Logf("> receiving passwords ...")
	for ok {
		select {
		case <-ctx.Done():
			assert.Fail(t, ctx.Err().Error())
			ok = false
		case pw, ok = <-ch:
			if ok {
				t.Logf("> ++ received %#v", pw)
				passwords = append(passwords, pw)
			}
		}
	}
	t.Logf("> receiving passwords ... done!")

	// verify received passwords
	expectedPasswords := []string{
		"AAA",
		"AAB",
		"ABA",
		"ABB",
		"BAA",
		"BAB",
		"BBA",
		"BBB",
	}
	assert.Equal(t, expectedPasswords, passwords)
}

func TestSequencer_Stream_Limited(t *testing.T) {
	s, err := NewSequencer(
		WithCharset("AB"),
		WithLength(3),
	)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ch := make(chan string, 1)
	// streams all possible passwords into the channel in an async routine
	go func() {
		t.Logf("< streaming passwords ...")
		err2 := s.Stream(ctx, ch, big.NewInt(5))
		assert.Nil(t, err2)
		t.Logf("< streaming passwords ... done!")
	}()

	// listen on channel for passwords until channel is closed, or until timeout
	pw, ok := "", true
	var passwords []string
	t.Logf("> receiving passwords ...")
	for ok {
		select {
		case <-ctx.Done():
			assert.Fail(t, ctx.Err().Error())
			ok = false
		case pw, ok = <-ch:
			if ok {
				t.Logf("> ++ received %#v", pw)
				passwords = append(passwords, pw)
			}
		}
	}
	t.Logf("> receiving passwords ... done!")

	// verify received passwords
	expectedPasswords := []string{
		"AAA",
		"AAB",
		"ABA",
		"ABB",
		"BAA",
	}
	assert.Equal(t, expectedPasswords, passwords)
}
