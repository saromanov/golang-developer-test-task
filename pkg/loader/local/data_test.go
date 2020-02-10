package local

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadData(t *testing.T) {
	l, err := New("../../../assets/data.json")
	assert.NoError(t, err)
	data, err := l.Load()
	assert.NoError(t, err)
	assert.NotEqual(t, len(data), 0)
}
