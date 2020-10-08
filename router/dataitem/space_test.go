package dataitem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSpaceDataItem(t *testing.T) {
	_, err := NewSpaceDataItem("wrong key", 15)
	assert.Error(t, err)

	got, err := NewSpaceDataItem("{\"Lon\":21.21, \"Lat\":42.42}", 15)
	assert.NoError(t, err)
	exp := SpaceDataItem{
		Key:  "{\"Lon\":21.21, \"Lat\":42.42}",
		Lat:  42.42,
		Lon:  21.21,
		size: 15,
	}
	assert.Equal(t, exp, got)
	assert.Equal(t, "{\"Lon\":21.21, \"Lat\":42.42}", got.ID())
	assert.Equal(t, []interface{}{42.42, 21.21}, got.Values())
	assert.Equal(t, 15, int(got.Size()))
}
