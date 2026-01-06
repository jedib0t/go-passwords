package passphrase

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	passphrase, err := Generate()
	assert.NoError(t, err)
	assert.NotEmpty(t, passphrase)

	// Verify structure: should have 3 words separated by "-"
	words := strings.Split(passphrase, "-")
	assert.Equal(t, 3, len(words), "passphrase should have 3 words: %s", passphrase)
}
