package password

import "errors"

var (
	ErrEmptyCharset       = errors.New("cannot generate passwords with empty charset")
	ErrInvalidN           = errors.New("value of N exceeds valid range")
	ErrMinSymbolsTooLong  = errors.New("minimum number of symbols requested longer than password")
	ErrNoSymbolsInCharset = errors.New("found no symbols to use in charset")
	ErrZeroLenPassword    = errors.New("cannot generate passwords with 0 length")
)
