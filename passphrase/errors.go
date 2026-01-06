package passphrase

import (
	"fmt"
)

var (
	ErrDictionaryTooSmall = fmt.Errorf("dictionary cannot have less than %d words after word-length restrictions are applied", MinWordsInDictionary)
	ErrNumWordsTooSmall   = fmt.Errorf("number of words cannot be less than %d", NumWordsMin)
	ErrNumWordsTooLarge   = fmt.Errorf("number of words cannot be more than %d", NumWordsMax)
	ErrWordLengthInvalid  = fmt.Errorf("word-length rule invalid")
)
