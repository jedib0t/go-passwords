package enumerator

import (
	"math/big"
	"sync"

	"github.com/jedib0t/go-passwords/charset"
)

var (
	biZero = big.NewInt(0)
	biOne  = big.NewInt(1)
)

// Enumerator defines interfaces to manipulate an Enumerator.
type Enumerator interface {
	// AtEnd returns true if the Enumerator is at the last possible value.
	AtEnd() bool
	// Decrement moves the gears back by one turn.
	Decrement() bool
	// DecrementN moves the gears back by N turns.
	DecrementN(n *big.Int) bool
	// First moves the gears to the first possible value.
	First()
	// GoTo moves the gears to a specific location in the list of possible
	// locations. The value of 'n' is 1-indexed.
	GoTo(n *big.Int) error
	// Increment moves the gears forward by one turn.
	Increment() bool
	// IncrementN moves the gears forward by N turns.
	IncrementN(n *big.Int) bool
	// Last moves the gears to the last possible value.
	Last()
	// Location returns the current location in the list of possible locations
	// and is 1-indexed.
	Location() *big.Int
	// String returns the value as an end-user would see when they look at an
	// Enumerator.
	String() string
}

type enumerator struct {
	base           int
	baseBigInt     *big.Int
	charset        []rune
	length         int
	location       *big.Int
	locationMax    *big.Int
	mutex          sync.RWMutex
	rollover       bool
	value          []int
	valueInCharset []rune
}

// New returns a new Enumerator with "length" gears each containing the given
// Charset as the values.
func New(cs charset.Charset, length int, opts ...Option) Enumerator {
	base := len(cs)
	maxValues := numValues(base, length)

	o := &enumerator{
		base:           base,
		baseBigInt:     big.NewInt(int64(base)),
		charset:        []rune(cs),
		length:         length,
		location:       big.NewInt(1),
		locationMax:    new(big.Int).Set(maxValues),
		value:          make([]int, length),
		valueInCharset: make([]rune, length),
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *enumerator) AtEnd() bool {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	return o.location.Cmp(o.locationMax) >= 0
}

func (o *enumerator) Decrement() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// set the location
	if o.location.Cmp(biOne) == 0 { // at first
		if o.rollover {
			o.last()
			return true
		}
		return false
	}

	// decrement value
	o.location.Sub(o.location, biOne)
	for idx := o.length - 1; idx >= 0; idx-- {
		if o.decrementAtIndex(idx) {
			return true
		}
	}
	return true
}

func (o *enumerator) DecrementN(n *big.Int) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.location.Sub(o.location, n)
	if o.location.Cmp(biOne) < 0 { // less than min
		if !o.rollover {
			o.first()
			return false
		}
		// move backwards from max; o.location is currently -ve --> so Add()
		for o.location.Cmp(biOne) < 0 {
			o.location.Add(o.locationMax, o.location)
		}
	}
	o.computeValue()
	return true
}

func (o *enumerator) First() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.first()
}

func (o *enumerator) Location() *big.Int {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	return new(big.Int).Set(o.location)
}

func (o *enumerator) Increment() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// set the location
	if o.location.Cmp(o.locationMax) == 0 { // at max
		if o.rollover {
			o.first()
			return true
		}
		return false
	}

	// increment value
	o.location.Add(o.location, biOne)
	for idx := o.length - 1; idx >= 0; idx-- {
		if o.incrementAtIndex(idx) {
			return true
		}
	}
	return true
}

func (o *enumerator) IncrementN(n *big.Int) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.location.Add(o.location, n)
	if o.location.Cmp(o.locationMax) > 0 { // more than max
		if !o.rollover {
			o.last()
			return false
		}
		// move forwards from zero
		for o.location.Cmp(o.locationMax) > 0 {
			o.location.Sub(o.location, o.locationMax)
		}
	}
	o.computeValue()
	return true
}

func (o *enumerator) Last() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.last()
}

func (o *enumerator) GoTo(n *big.Int) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if n.Cmp(biOne) < 0 || n.Cmp(o.locationMax) > 0 {
		return ErrInvalidLocation
	}
	o.location.Set(n)
	o.computeValue()
	return nil
}

func (o *enumerator) String() string {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	for idx := range o.valueInCharset {
		o.valueInCharset[idx] = o.charset[o.value[idx]]
	}
	return string(o.valueInCharset)
}

func (o *enumerator) computeValue() {
	// base conversion: convert the value of location to a decimal with the
	// given base using continuous division and use all the remainders as the
	// values

	// prep the dividend, remainder and modulus
	dividend, remainder := new(big.Int).Sub(o.location, biOne), new(big.Int)
	// append values in reverse (from right to left)
	valIdx := o.length - 1
	// append every remainder until dividend becomes zero
	for ; dividend.Cmp(biZero) != 0; valIdx-- {
		dividend, remainder = dividend.QuoRem(dividend, o.baseBigInt, remainder)
		o.value[valIdx] = int(remainder.Int64())
	}
	// left-pad the remaining characters with 0 (==> 0th char in charset)
	for ; valIdx >= 0; valIdx-- {
		o.value[valIdx] = 0
	}
}

func (o *enumerator) decrementAtIndex(idx int) bool {
	if o.value[idx] > 0 {
		o.value[idx]--
		return true
	}
	if o.value[idx] == 0 && idx > 0 {
		o.value[idx] = o.base - 1
		o.decrementAtIndex(idx - 1)
		return true
	}
	return false
}

func (o *enumerator) first() {
	o.location.Set(biOne)
	for idx := range o.value {
		o.value[idx] = 0
	}
}

func (o *enumerator) incrementAtIndex(idx int) bool {
	if o.value[idx] < o.base-1 {
		o.value[idx]++
		return true
	}
	if o.value[idx] == o.base-1 && idx > 0 {
		o.value[idx] = 0
		o.incrementAtIndex(idx - 1)
		return true
	}
	return false
}

func (o *enumerator) last() {
	o.location.Set(o.locationMax)
	for idx := range o.value {
		o.value[idx] = o.base - 1
	}
}
