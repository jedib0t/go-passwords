package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	assert.NotEmpty(t, Generate())

	SetSeed(1)
	assert.Equal(t, "rk&nkRHeg54P", Generate())
}
