package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"math/rand/v2"
	"sync"

	"github.com/jedib0t/go-passwords/charset"
)

var (
	biZero = big.NewInt(0)
	biOne  = big.NewInt(1)
)

// Sequencer is a deterministic Password Generator that generates all possible
// combinations of passwords for a Charset and defined number of characters in
// the password. It lets you move back and forth through the list of possible
// passwords, and involves no RNG.
type Sequencer interface {
	// First moves to the first possible password and returns the same.
	First() string
	// Get returns the current password in the sequence.
	Get() string
	// GetN returns the value for N (location in list of possible passwords).
	GetN() *big.Int
	// GotoN overrides N.
	GotoN(n *big.Int) (string, error)
	// HasNext returns true if there is at least one more password.
	HasNext() bool
	// Last moves to the last possible password and returns the same.
	Last() string
	// Next moves to the next possible password and returns the same.
	Next() string
	// NextN is like Next looped N times, in an optimal way.
	NextN(n *big.Int) string
	// Prev moves to the previous possible password and returns the same.
	Prev() string
	// PrevN is like Prev looped N times, in an optimal way.
	PrevN(n *big.Int) string
	// Reset cleans up state and moves to the first possible word.
	Reset()
	// Stream sends all possible passwords in order to the given channel. If you
	// want to limit output, pass in a *big.Int with the number of passwords you
	// want to be generated and streamed.
	Stream(ctx context.Context, ch chan string, optionalCount ...*big.Int) error
}

type sequencer struct {
	base           *big.Int
	charset        []rune
	charsetLen     int
	charsetMaxIdx  int
	maxWords       *big.Int
	mutex          sync.Mutex
	n              *big.Int
	nMax           *big.Int
	numChars       int
	password       []int
	passwordChars  []rune
	passwordMaxIdx int
	rng            *rand.Rand
}

// New returns a password Sequencer.
func New(rules ...Rule) (Sequencer, error) {
	s := &sequencer{}
	for _, rule := range append(basicRules, rules...) {
		rule(s)
	}

	// init the variables
	s.base = big.NewInt(int64(len(s.charset)))
	s.charsetLen = len(s.charset)
	s.charsetMaxIdx = len(s.charset) - 1
	s.maxWords = MaximumPossibleWords(charset.Charset(s.charset), s.numChars)
	s.n = big.NewInt(0)
	s.nMax = new(big.Int).Sub(s.maxWords, biOne)
	s.password = make([]int, s.numChars)
	s.passwordChars = make([]rune, s.numChars)
	s.passwordMaxIdx = s.numChars - 1

	if len(s.charset) == 0 {
		return nil, ErrEmptyCharset
	}
	if s.numChars <= 0 {
		return nil, ErrZeroLenPassword
	}
	return s, nil
}

// First moves to the first possible password and returns the same.
func (s *sequencer) First() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.n.Set(biZero)
	for idx := range s.password {
		s.password[idx] = 0
	}
	return s.get()
}

// Get returns the current password in the sequence.
func (s *sequencer) Get() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.get()
}

// GetN returns the current location in the list of possible passwords.
func (s *sequencer) GetN() *big.Int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return new(big.Int).Set(s.n)
}

// GotoN overrides the current location in the list of possible passwords.
func (s *sequencer) GotoN(n *big.Int) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// ensure n is within possible range (0 to nMax)
	if n.Sign() < 0 || n.Cmp(s.nMax) > 0 {
		return "", fmt.Errorf("%w: n=%s, range=[0 to %s]", ErrInvalidN, n, s.nMax)
	}

	// override and compute the word
	s.n.Set(n)
	s.computeWord()
	return s.get(), nil
}

// HasNext returns true if there is at least one more password.
func (s *sequencer) HasNext() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.n.Cmp(s.nMax) < 0
}

// Last moves to the last possible password and returns the same.
func (s *sequencer) Last() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.n.Set(s.nMax)
	for idx := range s.password {
		s.password[idx] = s.charsetMaxIdx
	}
	return s.get()
}

// Next moves to the next possible password and returns the same.
func (s *sequencer) Next() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.next()
	return s.get()
}

// NextN is like Next looped N times, in an optimal way.
func (s *sequencer) NextN(n *big.Int) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.n.Cmp(s.nMax) < 0 {
		s.n.Add(s.n, n)
		s.computeWord()
	}
	return s.get()
}

// Prev moves to the previous possible password and returns the same.
func (s *sequencer) Prev() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.prev()
	return s.get()
}

// PrevN is like Prev looped N times, in an optimal way.
func (s *sequencer) PrevN(n *big.Int) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.n.Cmp(biZero) > 0 {
		s.n.Sub(s.n, n)
		s.computeWord()
	}
	return s.get()
}

// Reset cleans up state and moves to the first possible word.
func (s *sequencer) Reset() {
	s.First()
}

// Stream sends all possible passwords in order to the given channel. If you
// want to limit output, pass in a *big.Int with the number of passwords you
// want to be generated and streamed.
func (s *sequencer) Stream(ctx context.Context, ch chan string, optionalCount ...*big.Int) error {
	defer close(ch)

	maxToBeSent := new(big.Int).Set(s.maxWords)
	if len(optionalCount) == 1 && optionalCount[0] != nil && optionalCount[0].Cmp(biZero) > 0 {
		maxToBeSent.Set(optionalCount[0])
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	ch <- s.get()
	chSent := big.NewInt(1)
	for ; s.next() && chSent.Cmp(maxToBeSent) < 0; chSent.Add(chSent, biOne) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ch <- s.get()
		}
	}
	return nil
}

func (s *sequencer) computeWord() {
	// base conversion: convert the value of n to a decimal with a base of
	// wg.charsetLen using continuous division and use all the remainders as the
	// index value to pick the character from the charset

	// prep the dividend, remainder and modulus
	dividend, remainder := new(big.Int).Set(s.n), new(big.Int)
	// append values to the password in reverse (from right to left)
	charIdx := s.passwordMaxIdx
	// append every remainder until dividend becomes zero
	for ; dividend.Cmp(biZero) != 0; charIdx-- {
		dividend, remainder = dividend.QuoRem(dividend, s.base, remainder)
		s.password[charIdx] = int(remainder.Int64())
	}
	// left-pad the remaining characters with 0 (==> 0th char in charset)
	for ; charIdx >= 0; charIdx-- {
		s.password[charIdx] = 0
	}
}

func (s *sequencer) get() string {
	for idx := range s.passwordChars {
		s.passwordChars[idx] = s.charset[s.password[idx]]
	}
	return string(s.passwordChars)
}

func (s *sequencer) next() bool {
	if s.n.Cmp(s.nMax) >= 0 {
		return false
	}

	s.n.Add(s.n, biOne)
	for idx := s.passwordMaxIdx; idx >= 0; idx-- {
		if s.nextAtIndex(idx) {
			return true
		}
	}
	return true
}

func (s *sequencer) nextAtIndex(idx int) bool {
	if s.password[idx] < s.charsetMaxIdx {
		s.password[idx]++
		return true
	}
	if s.password[idx] == s.charsetMaxIdx && idx > 0 {
		s.password[idx] = 0
		s.nextAtIndex(idx - 1)
		return true
	}
	return false
}

func (s *sequencer) prev() bool {
	if s.n.Cmp(biZero) <= 0 {
		return false
	}

	s.n.Sub(s.n, biOne)
	for idx := s.passwordMaxIdx; idx >= 0; idx-- {
		if s.prevAtIndex(idx) {
			return true
		}
	}
	return true
}

func (s *sequencer) prevAtIndex(idx int) bool {
	if s.password[idx] > 0 {
		s.password[idx]--
		return true
	}
	if s.password[idx] == 0 && idx > 0 {
		s.password[idx] = s.charsetMaxIdx
		s.prevAtIndex(idx - 1)
		return true
	}
	return false
}

func (s *sequencer) ruleEnforcer() {}
