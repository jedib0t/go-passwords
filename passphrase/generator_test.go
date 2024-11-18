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
		"Eaters-Likers6-Coiler",
		"Shags-Obeyed2-Tweeze",
		"Reared-Campos-Umbral5",
		"Subers-Bemean-Sall6",
		"Priory1-Prayer-Mirk",
		"Kills8-Alkoxy-Sequel",
		"Long-Reply-Coco4",
		"Embank3-Tusche-Degage",
		"Pitons-Luce9-Jabber",
		"Flavin-Capful-Leaved2",
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
