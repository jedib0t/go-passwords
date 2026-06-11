package enumerator

import (
	"math/big"
	"testing"

	"github.com/jedib0t/go-passwords/charset"
)

// FuzzEnumeratorOps drives an Enumerator through arbitrary operation
// sequences and asserts the core invariants after every step: the location
// stays within [1, base^length] and String always has exactly "length" runes
// drawn from the charset.
func FuzzEnumeratorOps(f *testing.F) {
	f.Add(uint8(0), uint8(2), []byte{0, 1, 2, 3, 4, 5, 6}, false)
	f.Add(uint8(1), uint8(3), []byte{5, 2, 2, 3, 0, 6, 1}, true)
	f.Fuzz(func(t *testing.T, csChoice uint8, length uint8, ops []byte, rollover bool) {
		charsets := []charset.Charset{
			charset.Numbers,
			charset.AlphabetsLower,
			charset.AlphaNumeric,
			charset.Charset("ab"),
		}
		cs := charsets[int(csChoice)%len(charsets)]
		l := int(length%6) + 1
		if len(ops) > 256 {
			t.Skip()
		}

		var opts []Option
		if rollover {
			opts = append(opts, WithRolloverEnabled(true))
		}
		e, err := New(cs, l, opts...)
		if err != nil {
			t.Fatalf("New(%q, %d) failed: %v", cs, l, err)
		}

		base := int64(len([]rune(cs)))
		max := new(big.Int).Exp(big.NewInt(base), big.NewInt(int64(l)), nil)
		for _, op := range ops {
			applyOp(e, op, max)
			if loc := e.Location(); loc.Sign() < 1 || loc.Cmp(max) > 0 {
				t.Fatalf("location %s out of range [1, %s] after op %d", loc, max, op)
			}
			s := e.String()
			runes := []rune(s)
			if len(runes) != l {
				t.Fatalf("String() = %q has %d runes, want %d", s, len(runes), l)
			}
			for _, r := range runes {
				if !cs.Contains(r) {
					t.Fatalf("String() = %q contains %q which is not in charset %q", s, r, cs)
				}
			}
		}
	})
}

func applyOp(e Enumerator, op byte, max *big.Int) {
	n := big.NewInt(int64(op))
	switch op % 7 {
	case 0:
		e.Increment()
	case 1:
		e.Decrement()
	case 2:
		e.IncrementN(n)
	case 3:
		e.DecrementN(n)
	case 4:
		e.First()
	case 5:
		e.Last()
	case 6:
		// GoTo with an arbitrary target; errors for out-of-range targets are
		// expected and fine, the invariants are checked by the caller
		_ = e.GoTo(n)
	}
}
