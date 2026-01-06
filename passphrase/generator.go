package passphrase

import (
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/jedib0t/go-passwords/rng"
)

const (
	MinWordsInDictionary = 256
	NumWordsMin          = 2
	NumWordsMax          = 32
)

type Generator interface {
	// Generate returns a randomly generated password.
	Generate() (string, error)
	// GenerateTo generates a password and writes it to the provided buffer.
	// It returns the number of bytes written or an error.
	GenerateTo(buf []byte) (int, error)
}

type generator struct {
	capitalize    bool
	dictionary    []string
	dictionaryLen int
	separator     string
	numWords      int
	withNumber    bool
	wordLenMin    int
	wordLenMax    int
}

// NewGenerator returns a password generator that implements the Generator
// interface.
func NewGenerator(rules ...Rule) (Generator, error) {
	g := &generator{}
	for _, opt := range append(basicRules, rules...) {
		opt(g)
	}
	return g.sanitize()
}

// Generate returns a randomly generated password.
func (g *generator) Generate() (string, error) {
	// estimate the capacity needed: max word length * num words + separators + digit
	buf := make([]byte, g.wordLenMax*g.numWords+len(g.separator)*(g.numWords-1)+1)
	n, err := g.GenerateTo(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

// GenerateTo generates a password and writes it to the provided buffer.
// It returns the number of bytes written or an error.
func (g *generator) GenerateTo(buf []byte) (int, error) {
	// inject a random number after one of the words if asked for
	wordForDigitSuffixIdx, digit := -1, 0
	if g.withNumber {
		var err error
		if wordForDigitSuffixIdx, err = rng.IntN(g.numWords); err != nil {
			return 0, err
		}
		if digit, err = rng.IntN(10); err != nil {
			return 0, err
		}
	}

	// Select unique word indices using rejection sampling
	var wordIndices [NumWordsMax]int
	for i := 0; i < g.numWords; i++ {
		wordIndex, err := g.getUniqueWordIndex(wordIndices[:i])
		if err != nil {
			return 0, err
		}
		wordIndices[i] = wordIndex
	}

	// append words to the buffer
	offset := 0
	for idx := 0; idx < g.numWords; idx++ {
		err := g.writeWordToBuf(buf, &offset, g.dictionary[wordIndices[idx]],
			idx == wordForDigitSuffixIdx, digit, idx < g.numWords-1)
		if err != nil {
			return 0, err
		}
	}

	return offset, nil
}

func (g *generator) getUniqueWordIndex(pickedIndices []int) (int, error) {
	for {
		wordIndex, err := rng.IntN(g.dictionaryLen)
		if err != nil {
			return 0, err
		}
		// check if already picked
		picked := false
		for _, pickedIndex := range pickedIndices {
			if pickedIndex == wordIndex {
				picked = true
				break
			}
		}
		if !picked {
			return wordIndex, nil
		}
	}
}

func (g *generator) writeWordToBuf(buf []byte, offset *int, word string, addDigit bool, digit int, addSeparator bool) error {
	if *offset+len(word) > len(buf) {
		return ErrBufferTooSmall
	}
	*offset += copy(buf[*offset:], word)

	if addDigit {
		if *offset+1 > len(buf) {
			return ErrBufferTooSmall
		}
		buf[*offset] = '0' + byte(digit)
		(*offset)++
	}

	if addSeparator {
		if *offset+len(g.separator) > len(buf) {
			return ErrBufferTooSmall
		}
		*offset += copy(buf[*offset:], g.separator)
	}

	return nil
}

func (g *generator) sanitize() (Generator, error) {
	// check if the word length is valid
	if g.wordLenMin < 1 || g.wordLenMin > g.wordLenMax {
		return nil, ErrWordLengthInvalid
	}

	// remove words that are too-short & too-long
	g.dictionary = slices.DeleteFunc(g.dictionary, func(word string) bool {
		return len(word) < g.wordLenMin || len(word) > g.wordLenMax
	})
	slices.Sort(g.dictionary)
	g.dictionary = slices.Compact(g.dictionary)
	g.dictionaryLen = len(g.dictionary)

	// check if the dictionary is too small
	if g.dictionaryLen < g.numWords || g.dictionaryLen < MinWordsInDictionary {
		return nil, ErrDictionaryTooSmall
	}

	// capitalize all words in the dictionary ahead of time
	if g.capitalize {
		for idx := range g.dictionary {
			r, size := utf8.DecodeRuneInString(g.dictionary[idx])
			if r != utf8.RuneError {
				g.dictionary[idx] = string(unicode.ToUpper(r)) + g.dictionary[idx][size:]
			}
		}
	}

	// check if the number of words is too small or too large
	if g.numWords < NumWordsMin {
		return nil, ErrNumWordsTooSmall
	}
	if g.numWords > NumWordsMax {
		return nil, ErrNumWordsTooLarge
	}
	return g, nil
}
