package cudainfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDevice ...
func TestDevice(t *testing.T) {
	v, err := GetCUDAVersion()
	assert.NoError(t, err)
	assert.NotEqual(t, "", v)
	t.Log(v)
}
