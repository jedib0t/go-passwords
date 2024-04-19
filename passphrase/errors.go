package passphrase

import (
	"errors"
	"fmt"
)

var (
	ErrDictionaryTooSmall = errors.New(fmt.Sprintf("dictionary should have more than %d words", MinWords))
	ErrNumWordsInvalid    = errors.New("number of words cannot be less than 1")
	ErrWordLengthInvalid  = errors.New("word-length rule invalid")
)
