package dictionaries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnglish(t *testing.T) {
	assert.NotEmpty(t, English())
}
