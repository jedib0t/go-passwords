package password

import "errors"

var (
	ErrEmptyCharset    = errors.New("cannot generate passwords with empty charset")
	ErrZeroLenPassword = errors.New("cannot generate passwords with 0 length")
	ErrInvalidN        = errors.New("value of N exceeds valid range")
)
