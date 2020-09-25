package transform

import (
	"github.com/stretchr/testify/assert"
	sfcmocks "github.com/struckoff/sfcframework/mocks"
	"testing"
)

func TestSpaceTransform(t *testing.T) {
	type args struct {
		values  []interface{}
		dimsize uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				values:  []interface{}{42.3, 121.21},
				dimsize: 4,
			},
			want:    []uint64{2, 3},
			wantErr: false,
		},
		{
			name: "lon limit",
			args: args{
				values:  []interface{}{100.11, 21.21},
				dimsize: 4,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "lat limit",
			args: args{
				values:  []interface{}{42.3, 221.21},
				dimsize: 4,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfc := &sfcmocks.Curve{}
			sfc.On("Dimensions").Return(uint64(2))
			sfc.On("DimensionSize").Return(tt.args.dimsize)
			got, err := SpaceTransform(tt.args.values, sfc)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpaceTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
