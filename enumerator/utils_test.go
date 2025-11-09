package enumerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_numValues(t *testing.T) {
	assert.Equal(t, "0", numValues(10, 0).String())
	assert.Equal(t, "10", numValues(10, 1).String())
	assert.Equal(t, "10000", numValues(10, 4).String())
	assert.Equal(t, "100000000", numValues(10, 8).String())
	assert.Equal(t, "10000000000", numValues(10, 10).String())
}
