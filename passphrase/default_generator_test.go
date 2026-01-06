package passphrase

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	passphrase := Generate()
	assert.NotEmpty(t, passphrase)

	// Verify structure: should have 3 words separated by "-"
	words := strings.Split(passphrase, "-")
	assert.Equal(t, 3, len(words), "passphrase should have 3 words: %s", passphrase)
}
