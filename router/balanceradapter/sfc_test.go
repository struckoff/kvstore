package balanceradapter

import (
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/optimizer"
	"github.com/struckoff/kvstore/router/transform"
	balancer "github.com/struckoff/sfcframework"
	"github.com/struckoff/sfcframework/curve"
	"github.com/struckoff/sfcframework/node"
	"testing"
)

func TestNewSFCBalancer(t *testing.T) {
	type args struct {
		conf *config.BalancerConfig
	}
	type bal struct {
		cType      curve.CurveType
		dims, size uint64
		tf         balancer.TransformFunc
		op         balancer.OptimizerFunc
		nodes      []node.Node
	}
	tests := []struct {
		name    string
		args    args
		want    bal
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				conf: &config.BalancerConfig{
					Mode: config.SFCMode,
					SFC: &config.SFCConfig{
						Dimensions: 2,
						Size:       4,
						Curve: config.CurveType{
							CurveType: curve.Morton,
						},
					},
					NodeHash: 1,
					DataMode: config.GeoData,
					State:    false,
				},
			},
			want: bal{
				cType: curve.Morton,
				dims:  2,
				size:  4,
				tf:    transform.SpaceTransform,
				op:    optimizer.RangeOptimizer,
				nodes: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSFCBalancer(tt.args.conf)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				bal, err := balancer.NewBalancer(tt.want.cType, tt.want.dims, tt.want.size, tt.want.tf, tt.want.op, tt.want.nodes)
				if err != nil {
					t.Fatal(err)
				}
				exp := &SFC{bal}

				assert.NoError(t, err)
				assert.Equal(t, exp.bal, got.bal)
			}
		})
	}
}
