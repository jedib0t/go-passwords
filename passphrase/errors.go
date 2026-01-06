package passphrase

import (
	"fmt"
)

var (
	ErrBufferTooSmall     = fmt.Errorf("buffer is too small to hold the generated passphrase")
	ErrDictionaryTooSmall = fmt.Errorf("dictionary cannot have less than %d words after word-length restrictions are applied", MinWordsInDictionary)
	ErrNumWordsTooLarge   = fmt.Errorf("number of words cannot be more than %d", NumWordsMax)
	ErrNumWordsTooSmall   = fmt.Errorf("number of words cannot be less than %d", NumWordsMin)
	ErrWordLengthInvalid  = fmt.Errorf("word-length rule invalid")
)
