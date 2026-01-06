package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	password := Generate()
	assert.NotEmpty(t, password)
	assert.Equal(t, 12, len(password), "password should be 12 characters long")
}
