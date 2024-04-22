package passphrase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	assert.NotEmpty(t, Generate())

	SetSeed(1)
	assert.Equal(t, "Holism-Sharing6-Elates", Generate())
}
