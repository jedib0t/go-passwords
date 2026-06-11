package enumerator

import "errors"

var (
	ErrEmptyCharset    = errors.New("cannot enumerate with an empty charset")
	ErrInvalidLength   = errors.New("cannot enumerate with a length less than 1")
	ErrInvalidLocation = errors.New("invalid location")
)
