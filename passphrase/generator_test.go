package passphrase

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jedib0t/go-passwords/passphrase/dictionaries"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerator_Generate(b *testing.B) {
	g, err := NewGenerator()
	assert.Nil(b, err)
	assert.NotEmpty(b, g.Generate())

	for idx := 0; idx < b.N; idx++ {
		_ = g.Generate()
	}
}

func TestGenerator_Generate(t *testing.T) {
	g, err := NewGenerator(
		WithCapitalizedWords(true),
		WithDictionary(dictionaries.English()),
		WithNumWords(3),
		WithNumber(true),
		WithSeparator("-"),
		WithWordLength(4, 6),
	)
	assert.NotNil(t, g)
	assert.Nil(t, err)
	g.SetSeed(1)

	expectedPassphrases := []string{
		"Sans-Liber-Quale1",
		"Defogs-Tael0-Hallo",
		"Medium-Leader-Sesame2",
		"Chelae-Tocsin8-Haling",
		"Taxies1-Sordor-Banner",
		"Kwanza-Molies-Lapses5",
		"Scurf-Hookas-Beryl4",
		"Repine-Dele-Loans3",
		"Furore0-Geneva-Celts",
		"Strew7-Tweed-Sannop",
		"Quasi7-Vino-Optic",
		"Alible8-Sherds-Fraena",
	}
	var actualPhrases []string
	for idx := 0; idx < 1000; idx++ {
		passphrase := g.Generate()
		assert.NotEmpty(t, passphrase)
		if idx < len(expectedPassphrases) {
			actualPhrases = append(actualPhrases, passphrase)
			assert.Equal(t, expectedPassphrases[idx], passphrase)
		}
	}
	if !slices.Equal(expectedPassphrases, actualPhrases) {
		for _, pw := range actualPhrases {
			fmt.Printf("%#v,\n", pw)
		}
	}
}
