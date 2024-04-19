package passphrase

import "github.com/jedib0t/go-passwords/passphrase/dictionaries"

var (
	defaultGenerator, _ = NewGenerator(
		WithCapitalizedWords(true),
		WithDictionary(dictionaries.English()),
		WithNumWords(3),
		WithNumber(true),
		WithSeparator("-"),
		WithWordLength(4, 7),
	)
)

// Generate generates and returns a passphrase that follows the following rules:
// * uses an English Dictionary
// * uses capitalized words
// * uses a total of 3 words
// * injects a random number behind one of the words
// * uses "-" as the separator
// * ensures words used are between 4 and 7 characters long
func Generate() string {
	return defaultGenerator.Generate()
}

// SetSeed sets the seed value used by the RNG of the default Generator.
func SetSeed(seed uint64) {
	defaultGenerator.SetSeed(seed)
}
