package enumerator

import (
	"errors"
	"math/big"
	"testing"

	"github.com/jedib0t/go-passwords/charset"
	"github.com/stretchr/testify/assert"
)

func TestEnumerator(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		o := New(charset.Numbers, 2)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())

		ok := o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "2", o.Location().String())
		assert.Equal(t, "01", o.String())
		ok = o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "3", o.Location().String())
		assert.Equal(t, "02", o.String())

		o.Last()
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())
		ok = o.Increment()
		assert.False(t, ok)
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())

		ok = o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, "99", o.Location().String())
		assert.Equal(t, "98", o.String())
		ok = o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, "98", o.Location().String())
		assert.Equal(t, "97", o.String())

		o.First()
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())
		ok = o.Decrement()
		assert.False(t, ok)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())
	})

	t.Run("rollover", func(t *testing.T) {
		o := New(charset.Numbers, 2, WithRolloverEnabled(true))
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())

		ok := o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "2", o.Location().String())
		assert.Equal(t, "01", o.String())
		ok = o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "3", o.Location().String())
		assert.Equal(t, "02", o.String())

		o.Last()
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())
		ok = o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())
		ok = o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "2", o.Location().String())
		assert.Equal(t, "01", o.String())

		ok = o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())
		ok = o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())
		ok = o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, "99", o.Location().String())
		assert.Equal(t, "98", o.String())
	})

	t.Run("really big enumerator", func(t *testing.T) {
		o := New(charset.AllChars, 256)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", o.String())

		o.Last()
		assert.Equal(t, "****************************************************************************************************************************************************************************************************************************************************************", o.String())
		assert.Equal(t, "22135954000460481554501886154749459371625170502600730699163663905247049740079899968480034338379403807827944552623126075988673634259405600148560278663819464589512058373791164736632467335096807212642462431896323483136010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", o.Location().String())
		ok := o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, "***************************************************************************************************************************************************************************************************************************************************************&", o.String())
		assert.Equal(t, "22135954000460481554501886154749459371625170502600730699163663905247049740079899968480034338379403807827944552623126075988673634259405600148560278663819464589512058373791164736632467335096807212642462431896323483136009999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", o.Location().String())
		ok = o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, "***************************************************************************************************************************************************************************************************************************************************************^", o.String())
		assert.Equal(t, "22135954000460481554501886154749459371625170502600730699163663905247049740079899968480034338379403807827944552623126075988673634259405600148560278663819464589512058373791164736632467335096807212642462431896323483136009999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999998", o.Location().String())

		o.First()
		assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", o.String())
		assert.Equal(t, "1", o.Location().String())
		ok = o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB", o.String())
		assert.Equal(t, "2", o.Location().String())
		ok = o.Increment()
		assert.True(t, ok)
		assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC", o.String())
		assert.Equal(t, "3", o.Location().String())
	})
}

func TestEnumerator_Decrement(t *testing.T) {
	o := New(charset.Numbers, 3)
	assert.Equal(t, "1", o.Location().String())
	assert.Equal(t, "000", o.String())

	err := o.GoTo(big.NewInt(1000))
	assert.Nil(t, err)
	assert.Equal(t, "999", o.String())

	for idx := int64(999); idx >= 1; idx-- {
		ok := o.Decrement()
		assert.True(t, ok)
		assert.Equal(t, big.NewInt(idx), o.Location(), idx)
	}
}

func TestEnumerator_DecrementN(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		o := New(charset.Numbers, 2)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())

		ok := o.DecrementN(big.NewInt(5))
		assert.False(t, ok)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())

		o.Last()
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())

		ok = o.DecrementN(big.NewInt(5))
		assert.True(t, ok)
		assert.Equal(t, "95", o.Location().String())
		assert.Equal(t, "94", o.String())

		ok = o.DecrementN(big.NewInt(500))
		assert.False(t, ok)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())
	})

	t.Run("rollover", func(t *testing.T) {
		o := New(charset.Numbers, 2, WithRolloverEnabled(true))
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())

		ok := o.DecrementN(big.NewInt(1))
		assert.True(t, ok)
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())

		ok = o.DecrementN(big.NewInt(5))
		assert.True(t, ok)
		assert.Equal(t, "95", o.Location().String())
		assert.Equal(t, "94", o.String())

		ok = o.DecrementN(big.NewInt(500))
		assert.True(t, ok)
		assert.Equal(t, "95", o.Location().String())
		assert.Equal(t, "94", o.String())
	})
}

func TestEnumerator_GoTo(t *testing.T) {
	o := New(charset.Numbers, 2)
	assert.Equal(t, "1", o.Location().String())
	assert.Equal(t, "00", o.String())

	err := o.GoTo(big.NewInt(0))
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidLocation))

	err = o.GoTo(big.NewInt(5))
	assert.Nil(t, err)
	assert.Equal(t, "5", o.Location().String())
	assert.Equal(t, "04", o.String())

	err = o.GoTo(big.NewInt(50))
	assert.Nil(t, err)
	assert.Equal(t, "50", o.Location().String())
	assert.Equal(t, "49", o.String())

	err = o.GoTo(big.NewInt(100))
	assert.Nil(t, err)
	assert.Equal(t, "100", o.Location().String())
	assert.Equal(t, "99", o.String())

	err = o.GoTo(big.NewInt(101))
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidLocation))
}

func TestEnumerator_Increment(t *testing.T) {
	o := New(charset.Numbers, 3)
	assert.Equal(t, "1", o.Location().String())

	for idx := int64(2); idx <= 1000; idx++ {
		ok := o.Increment()
		assert.True(t, ok)
		assert.Equal(t, big.NewInt(idx), o.Location(), idx)
	}
}

func TestEnumerator_IncrementN(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		o := New(charset.Numbers, 2)
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())

		ok := o.IncrementN(big.NewInt(5))
		assert.True(t, ok)
		assert.Equal(t, "6", o.Location().String())
		assert.Equal(t, "05", o.String())

		o.Last()
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())

		ok = o.IncrementN(big.NewInt(5))
		assert.False(t, ok)
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())
	})

	t.Run("rollover", func(t *testing.T) {
		o := New(charset.Numbers, 2, WithRolloverEnabled(true))
		assert.Equal(t, "1", o.Location().String())
		assert.Equal(t, "00", o.String())

		ok := o.IncrementN(big.NewInt(5))
		assert.True(t, ok)
		assert.Equal(t, "6", o.Location().String())
		assert.Equal(t, "05", o.String())

		o.Last()
		assert.Equal(t, "100", o.Location().String())
		assert.Equal(t, "99", o.String())

		ok = o.IncrementN(big.NewInt(5))
		assert.True(t, ok)
		assert.Equal(t, "5", o.Location().String())
		assert.Equal(t, "04", o.String())
	})
}

func TestEnumerator_AtEnd(t *testing.T) {
	o := New(charset.Numbers, 2)
	assert.False(t, o.AtEnd())

	o.Last()
	assert.True(t, o.AtEnd())

	o.First()
	assert.False(t, o.AtEnd())

	o.Increment()
	assert.False(t, o.AtEnd())
}

func TestEnumerator_decrementAtIndex_EdgeCases(t *testing.T) {
	o := New(charset.Numbers, 3)

	// Test decrementAtIndex edge case by trying to decrement when already at first
	o.First()
	// Try to decrement - should fail since we're at first
	ok := o.Decrement()
	assert.False(t, ok)
	// Verify we're still at first
	assert.Equal(t, "000", o.String())
	assert.Equal(t, "1", o.Location().String())
}

func TestEnumerator_incrementAtIndex_EdgeCases(t *testing.T) {
	o := New(charset.Numbers, 3)

	// Test incrementAtIndex edge case by trying to increment when already at last
	o.Last()
	// Try to increment - should fail since we're at last
	ok := o.Increment()
	assert.False(t, ok)
	// Verify we're still at last
	maxLoc := o.Location()
	o.GoTo(big.NewInt(1000))
	assert.Equal(t, maxLoc.String(), o.Location().String())
}

func TestEnumerator_String_Cache(t *testing.T) {
	o := New(charset.Numbers, 2)

	// First call should compute and cache
	str1 := o.String()
	assert.Equal(t, "00", str1)

	// Second call should use cache
	str2 := o.String()
	assert.Equal(t, "00", str2)
	assert.Equal(t, str1, str2)

	// After increment, cache should be invalidated
	o.Increment()
	str3 := o.String()
	assert.Equal(t, "01", str3)
	assert.NotEqual(t, str1, str3)
}

func TestEnumerator_ensureLocation(t *testing.T) {
	o := New(charset.Numbers, 2)

	// Initially location should be clean
	loc1 := o.Location()
	assert.Equal(t, "1", loc1.String())

	// Increment should mark location as dirty
	o.Increment()
	// Location should be computed lazily
	loc2 := o.Location()
	assert.Equal(t, "2", loc2.String())

	// Multiple increments without calling Location
	o.Increment()
	o.Increment()
	o.Increment()
	// Location should still compute correctly
	loc3 := o.Location()
	assert.Equal(t, "5", loc3.String())
}

func TestEnumerator_Decrement_LoopCompletion(t *testing.T) {
	// Test case where decrementAtIndex doesn't return early in the loop
	// This covers the case where the loop completes without early return
	o := New(charset.Numbers, 1)
	o.GoTo(big.NewInt(5))

	// Decrement should work normally
	ok := o.Decrement()
	assert.True(t, ok)
	assert.Equal(t, "4", o.Location().String())
}

func TestEnumerator_Increment_LoopCompletion(t *testing.T) {
	// Test case where incrementAtIndex doesn't return early in the loop
	// This covers the case where the loop completes without early return
	o := New(charset.Numbers, 1)
	o.GoTo(big.NewInt(5))

	// Increment should work normally
	ok := o.Increment()
	assert.True(t, ok)
	assert.Equal(t, "6", o.Location().String())
}

func TestEnumerator_decrementAtIndex_RecursiveFailure(t *testing.T) {
	// Test the case where decrementAtIndex recursively calls itself but returns false
	// This happens when we try to decrement beyond the first position
	o := New(charset.Numbers, 2)
	o.First()

	// Try to decrement - should fail
	ok := o.Decrement()
	assert.False(t, ok)

	// Verify we're still at first
	assert.Equal(t, "00", o.String())
	assert.Equal(t, "1", o.Location().String())
}

func TestEnumerator_incrementAtIndex_RecursiveFailure(t *testing.T) {
	// Test the case where incrementAtIndex recursively calls itself but returns false
	// This happens when we try to increment beyond the last position
	o := New(charset.Numbers, 2)
	o.Last()

	// Try to increment - should fail
	ok := o.Increment()
	assert.False(t, ok)

	// Verify we're still at last
	assert.Equal(t, "99", o.String())
	maxLoc := o.Location()
	o.GoTo(big.NewInt(100))
	assert.Equal(t, maxLoc.String(), o.Location().String())
}

func TestEnumerator_Decrement_AllPaths(t *testing.T) {
	// Test Decrement when the loop doesn't find a match immediately
	// This covers the path where decrementAtIndex returns false and loop continues
	o := New(charset.Numbers, 3)
	o.GoTo(big.NewInt(100)) // "099"

	// Decrement should work
	ok := o.Decrement()
	assert.True(t, ok)
	assert.Equal(t, "098", o.String())
}

func TestEnumerator_Increment_AllPaths(t *testing.T) {
	// Test Increment when the loop doesn't find a match immediately
	// This covers the path where incrementAtIndex returns false and loop continues
	o := New(charset.Numbers, 3)
	o.GoTo(big.NewInt(100)) // "099"

	// Increment should work
	ok := o.Increment()
	assert.True(t, ok)
	assert.Equal(t, "100", o.String())
}
