package dataitem

import (
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/kvstore/router/config"
	"reflect"
	"testing"
)

func TestGetDataItemFunc(t *testing.T) {
	type args struct {
		dmt config.DataModeType
	}
	tests := []struct {
		name    string
		args    args
		want    NewDataItemFunc
		wantErr bool
	}{
		{
			name: "0",
			args: args{
				dmt: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "geo",
			args: args{
				dmt: config.GeoData,
			},
			want:    NewSpaceDataItem,
			wantErr: false,
		},
		{
			name: "kv",
			args: args{
				dmt: config.KVData,
			},
			want:    NewKVDataItem,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDataItemFunc(tt.args.dmt)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataItemFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotf := reflect.ValueOf(got)
			wantf := reflect.ValueOf(tt.want)
			assert.Equal(t, gotf, wantf)
		})
	}
}
