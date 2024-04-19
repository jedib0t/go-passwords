package passphrase

import (
	"errors"
	"fmt"
)

var (
	ErrDictionaryTooSmall = errors.New(fmt.Sprintf("dictionary cannot have less than %d words after word-length restrictions are applied", MinWordsInDictionary))
	ErrNumWordsInvalid    = errors.New(fmt.Sprintf("number of words cannot be less than %d", MinNumWords))
	ErrWordLengthInvalid  = errors.New("word-length rule invalid")
)
