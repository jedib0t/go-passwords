package password

import "github.com/jedib0t/go-passwords/charset"

var (
	defaultGenerator, _ = NewGenerator(
		WithCharset(charset.AllChars.WithoutAmbiguity()),
		WithLength(12),
		WithMinLowerCase(3),
		WithMinUpperCase(1),
		WithNumSymbols(1, 1),
	)
)

// Generate generates and returns a password that follows the following rules:
// * uses the AllChars charset
// * ensures the password is 12 characters long
// * uses a minimum of 3 lower case characters
// * uses a minimum of 1 upper case character
// * uses one symbol character
func Generate() string {
	return defaultGenerator.Generate()
}
