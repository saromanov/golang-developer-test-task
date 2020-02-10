package local

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadData(t *testing.T) {
	data, err := LocalLoad("../assets/data.json")
	assert.NoError(t, err)
	assert.NotEqual(t, len(data), 0)
}
