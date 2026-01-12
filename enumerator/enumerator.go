package enumerator

import (
	"math/big"
	"strings"
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
	base        int
	baseBigInt  *big.Int
	charset     []rune
	length      int
	location    *big.Int
	locationMax *big.Int
	mutex       sync.RWMutex
	rollover    bool
	value       []int
	isASCII     bool

	// Optimization: uint64 fast-path
	useUint64         bool
	locationUint64    uint64
	locationMaxUint64 uint64

	// Optimization: lazy location computation
	locationDirty bool

	// Optimization: reusable big.Int objects for computeValue and ensureLocation
	dividend   *big.Int
	remainder  *big.Int
	multiplier *big.Int
	val        *big.Int

	// Optimization: cached string result
	cachedString string
	stringDirty  bool
}

// New returns a new Enumerator with "length" gears each containing the given
// Charset as the values.
func New(cs charset.Charset, length int, opts ...Option) Enumerator {
	base := len(cs)
	maxValues := numValues(base, length)

	o := &enumerator{
		base:          base,
		baseBigInt:    big.NewInt(int64(base)),
		charset:       []rune(cs),
		length:        length,
		location:      big.NewInt(1),
		locationMax:   new(big.Int).Set(maxValues),
		value:         make([]int, length),
		dividend:      new(big.Int),
		remainder:     new(big.Int),
		multiplier:    new(big.Int),
		val:           new(big.Int),
		locationDirty: false,
		stringDirty:   true, // Need to compute initial string
	}

	// Detect if we can use uint64 fast-path
	if maxValues.IsUint64() {
		o.useUint64 = true
		o.locationUint64 = 1
		o.locationMaxUint64 = maxValues.Uint64()
	}

	// Check if charset is ASCII only for performance in String()
	o.isASCII = true
	for _, r := range o.charset {
		if r > 127 {
			o.isASCII = false
			break
		}
	}

	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *enumerator) AtEnd() bool {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	maxVal := o.base - 1
	for idx := range o.value {
		if o.value[idx] != maxVal {
			return false
		}
	}
	return true
}

func (o *enumerator) Decrement() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Check if at first by checking if all values are 0
	isFirst := true
	for idx := range o.value {
		if o.value[idx] != 0 {
			isFirst = false
			break
		}
	}
	if isFirst {
		if o.rollover {
			o.last()
			return true
		}
		return false
	}

	// Decrement value array directly
	for idx := o.length - 1; idx >= 0; idx-- {
		if o.decrementAtIndex(idx) {
			if !o.locationDirty {
				if o.useUint64 {
					o.locationUint64--
				}
				o.location.Sub(o.location, biOne)
			}
			o.stringDirty = true
			return true
		}
	}
	return true
}

func (o *enumerator) DecrementN(n *big.Int) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.ensureLocation()
	if o.useUint64 {
		nUint64 := n.Uint64()
		if o.locationUint64 > nUint64 {
			o.locationUint64 -= nUint64
		} else {
			if !o.rollover {
				o.first()
				return false
			}
			// rollover
			for o.locationUint64 <= nUint64 {
				o.locationUint64 += o.locationMaxUint64
			}
			o.locationUint64 -= nUint64
		}
		o.location.SetUint64(o.locationUint64)
	} else {
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
	}
	o.computeValue()
	o.locationDirty = false // location is now in sync with value
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

	o.ensureLocation()
	if o.useUint64 {
		return new(big.Int).SetUint64(o.locationUint64)
	}
	return new(big.Int).Set(o.location)
}

func (o *enumerator) Increment() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Check if at max by checking if all values are at max
	isMax := true
	maxVal := o.base - 1
	for idx := range o.value {
		if o.value[idx] != maxVal {
			isMax = false
			break
		}
	}
	if isMax {
		if o.rollover {
			o.first()
			return true
		}
		return false
	}

	// Increment value array directly
	for idx := o.length - 1; idx >= 0; idx-- {
		if o.incrementAtIndex(idx) {
			if !o.locationDirty {
				if o.useUint64 {
					o.locationUint64++
				}
				o.location.Add(o.location, biOne)
			}
			o.stringDirty = true
			return true
		}
	}
	return true
}

func (o *enumerator) IncrementN(n *big.Int) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.ensureLocation()
	if o.useUint64 {
		nUint64 := n.Uint64()
		o.locationUint64 += nUint64
		if o.locationUint64 > o.locationMaxUint64 {
			if !o.rollover {
				o.last()
				return false
			}
			// rollover
			for o.locationUint64 > o.locationMaxUint64 {
				o.locationUint64 -= o.locationMaxUint64
			}
		}
		o.location.SetUint64(o.locationUint64)
	} else {
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
	}
	o.computeValue()
	o.locationDirty = false // location is now in sync with value
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
	if o.useUint64 {
		o.locationUint64 = n.Uint64()
		o.location.SetUint64(o.locationUint64)
	} else {
		o.location.Set(n)
	}
	o.computeValue()
	o.locationDirty = false // location is now in sync with value
	return nil
}

func (o *enumerator) String() string {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.stringDirty {
		if o.isASCII {
			// Fast path for ASCII charsets
			b := make([]byte, o.length)
			for idx := range o.value {
				b[idx] = byte(o.charset[o.value[idx]])
			}
			o.cachedString = string(b)
		} else {
			var b strings.Builder
			b.Grow(o.length)
			for idx := range o.value {
				b.WriteRune(o.charset[o.value[idx]])
			}
			o.cachedString = b.String()
		}
		o.stringDirty = false
	}
	return o.cachedString
}

func (o *enumerator) computeValue() {
	// base conversion: convert the value of location to a decimal with the
	// given base using continuous division and use all the remainders as the
	// values

	if o.useUint64 {
		val := o.locationUint64 - 1
		valIdx := o.length - 1
		base := uint64(o.base)
		for ; val != 0; valIdx-- {
			o.value[valIdx] = int(val % base)
			val /= base
		}
		for ; valIdx >= 0; valIdx-- {
			o.value[valIdx] = 0
		}
	} else {
		// Reuse pre-allocated big.Int objects instead of creating new ones
		o.dividend.Sub(o.location, biOne)
		valIdx := o.length - 1
		// append every remainder until dividend becomes zero
		for ; o.dividend.Cmp(biZero) != 0; valIdx-- {
			o.dividend.QuoRem(o.dividend, o.baseBigInt, o.remainder)
			o.value[valIdx] = int(o.remainder.Int64())
		}
		// left-pad the remaining characters with 0 (==> 0th char in charset)
		for ; valIdx >= 0; valIdx-- {
			o.value[valIdx] = 0
		}
	}
	o.stringDirty = true
}

func (o *enumerator) decrementAtIndex(idx int) bool {
	if o.value[idx] > 0 {
		o.value[idx]--
		return true
	}
	if o.value[idx] == 0 && idx > 0 {
		o.value[idx] = o.base - 1
		if o.decrementAtIndex(idx - 1) {
			return true
		}
	}
	return false
}

func (o *enumerator) first() {
	if o.useUint64 {
		o.locationUint64 = 1
	}
	o.location.Set(biOne)
	for idx := range o.value {
		o.value[idx] = 0
	}
	o.locationDirty = false
	o.stringDirty = true
}

func (o *enumerator) incrementAtIndex(idx int) bool {
	if o.value[idx] < o.base-1 {
		o.value[idx]++
		return true
	}
	if o.value[idx] == o.base-1 && idx > 0 {
		o.value[idx] = 0
		if o.incrementAtIndex(idx - 1) {
			return true
		}
	}
	return false
}

func (o *enumerator) last() {
	if o.useUint64 {
		o.locationUint64 = o.locationMaxUint64
	}
	o.location.Set(o.locationMax)
	for idx := range o.value {
		o.value[idx] = o.base - 1
	}
	o.locationDirty = false
	o.stringDirty = true
}

// ensureLocation computes location from value array if it's dirty
func (o *enumerator) ensureLocation() {
	if !o.locationDirty {
		return
	}

	// Compute location from value array
	// location = 1 + sum(value[i] * base^(length-1-i))
	if o.useUint64 {
		o.locationUint64 = 1
		multiplier := uint64(1)
		base := uint64(o.base)
		for idx := o.length - 1; idx >= 0; idx-- {
			if o.value[idx] > 0 {
				o.locationUint64 += uint64(o.value[idx]) * multiplier
			}
			multiplier *= base
		}
		o.location.SetUint64(o.locationUint64)
	} else {
		o.location.Set(biOne)
		o.multiplier.SetInt64(1)
		for idx := o.length - 1; idx >= 0; idx-- {
			if o.value[idx] > 0 {
				o.val.SetInt64(int64(o.value[idx]))
				o.val.Mul(o.val, o.multiplier)
				o.location.Add(o.location, o.val)
			}
			o.multiplier.Mul(o.multiplier, o.baseBigInt)
		}
	}
	o.locationDirty = false
}
