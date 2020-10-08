package dataitem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKVDataItem(t *testing.T) {
	got, err := NewKVDataItem("test-di", 15)
	assert.NoError(t, err)

	exp := KVDataItem{"test-di", 15}
	assert.Equal(t, exp, got)
	assert.Equal(t, "test-di", got.ID())
	assert.Equal(t, 15, int(got.Size()))
	assert.Equal(t, []interface{}{"test-di"}, got.Values())
}
