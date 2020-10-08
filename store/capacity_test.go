package store

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/kvstore/router/rpcapi"
	"testing"
)

func TestCapacity_RPCGet(t *testing.T) {
	c := &Capacity{
		c: 42.42,
	}
	got, err := c.RPCGet(context.TODO(), &rpcapi.Empty{})
	assert.NoError(t, err)
	exp := &rpcapi.Capacity{Capacity: 42.42}
	assert.Equal(t, exp, got)
}
