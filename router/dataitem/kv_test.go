package dataitem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKVDataItem(t *testing.T) {
	got, err := NewKVDataItem("test-di")
	assert.NoError(t, err)

	exp := KVDataItem("test-di")
	assert.Equal(t, exp, got)
	assert.Equal(t, "test-di", got.ID())
	assert.Equal(t, 1, int(got.Size()))
	assert.Equal(t, []interface{}{"test-di"}, got.Values())
}
