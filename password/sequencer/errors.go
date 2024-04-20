package sequencer

import "errors"

var (
	ErrEmptyCharset    = errors.New("cannot generate passwords with empty charset")
	ErrInvalidN        = errors.New("value of N exceeds valid range")
	ErrZeroLenPassword = errors.New("cannot generate passwords with 0 length")
)
