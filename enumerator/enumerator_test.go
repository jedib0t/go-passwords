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
