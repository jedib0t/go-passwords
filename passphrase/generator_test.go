package passphrase

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jedib0t/go-passwords/passphrase/dictionaries"
	"github.com/stretchr/testify/assert"
)

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
		"Sannup-Libels-Quaky1",
		"Defog-Tads0-Hallel",
		"Medina-Leaden-Servos2",
		"Chela-Tocher8-Halids",
		"Taxied1-Sordid-Banned",
		"Kwacha-Molest-Lapser5",
		"Scups-Hookah-Berths4",
		"Repin-Delays-Loaner3",
		"Furor0-Genets-Celt",
		"Stress7-Twee-Sank",
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
