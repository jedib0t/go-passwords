package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	password, err := Generate()
	assert.NoError(t, err)
	assert.NotEmpty(t, password)
	assert.Equal(t, 12, len(password), "password should be 12 characters long")
}
