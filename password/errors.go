package password

import "errors"

var (
	ErrEmptyCharset         = errors.New("cannot generate passwords with empty charset")
	ErrInvalidN             = errors.New("value of N exceeds valid range")
	ErrMinLowerCaseTooLong  = errors.New("minimum number of lower-case characters requested longer than password")
	ErrMinUpperCaseTooLong  = errors.New("minimum number of upper-case characters requested longer than password")
	ErrMinSymbolsTooLong    = errors.New("minimum number of symbols requested longer than password")
	ErrNoLowerCaseInCharset = errors.New("found no lower-case characters to use in charset")
	ErrNoUpperCaseInCharset = errors.New("found no upper-case characters to use in charset")
	ErrNoSymbolsInCharset   = errors.New("found no symbols to use in charset")
	ErrRequirementsNotMet   = errors.New("minimum number of lower-case+upper-case+symbols requested longer than password")
	ErrZeroLenPassword      = errors.New("cannot generate passwords with 0 length")
)
