package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyncMap_Copy(t *testing.T) {
	type fields struct {
		s map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]string
	}{
		{
			name: "test",
			fields: fields{
				s: map[string][]string{
					"key-0": {"val-0-0", "val-0-1"},
					"key-1": {"val-1-0", "val-1-1"},
				},
			},
			want: map[string][]string{
				"key-0": {"val-0-0", "val-0-1"},
				"key-1": {"val-1-0", "val-1-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &SyncMap{
				s: tt.fields.s,
			}
			got := sm.Copy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSyncMap_Put(t *testing.T) {
	type fields struct {
		s map[string][]string
	}
	type args struct {
		key   string
		value []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string][]string
	}{
		{
			name: "create",
			fields: fields{
				s: map[string][]string{
					"key-0": {"val-0-0", "val-0-1"},
				},
			},
			args: args{
				key:   "key-1",
				value: []string{"val-1-0", "val-1-1"},
			},
			want: map[string][]string{
				"key-0": {"val-0-0", "val-0-1"},
				"key-1": {"val-1-0", "val-1-1"},
			},
		},
		{
			name: "nil val",
			fields: fields{
				s: map[string][]string{
					"key-0": {"val-0-0", "val-0-1"},
				},
			},
			args: args{
				key:   "key-1",
				value: nil,
			},
			want: map[string][]string{
				"key-0": {"val-0-0", "val-0-1"},
				"key-1": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &SyncMap{
				s: tt.fields.s,
			}

			sm.Put(tt.args.key, tt.args.value)

			assert.Equal(t, tt.want, sm.s)
		})
	}
}

func TestSyncMap_Get(t *testing.T) {
	type fields struct {
		s map[string][]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
		wantOk bool
	}{
		{
			name: "found",
			fields: fields{
				s: map[string][]string{
					"key-0": {"val-0-0", "val-0-1"},
					"key-1": {"val-1-0", "val-1-1"},
				},
			},
			args: args{
				key: "key-0",
			},
			want:   []string{"val-0-0", "val-0-1"},
			wantOk: true,
		},
		{
			name: "not found",
			fields: fields{
				s: map[string][]string{
					"key-0": {"val-0-0", "val-0-1"},
					"key-1": {"val-1-0", "val-1-1"},
				},
			},
			args: args{
				key: "key-not-found",
			},
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &SyncMap{
				s: tt.fields.s,
			}
			got, ok := sm.Get(tt.args.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestSyncMap_Delete(t *testing.T) {
	type fields struct {
		s map[string][]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string][]string
	}{
		{
			name: "create",
			fields: fields{
				s: map[string][]string{
					"key-0": {"val-0-0", "val-0-1"},
					"key-1": {"val-1-0", "val-1-1"},
				},
			},
			args: args{
				key: "key-1",
			},
			want: map[string][]string{
				"key-0": {"val-0-0", "val-0-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &SyncMap{
				s: tt.fields.s,
			}

			sm.Delete(tt.args.key)

			assert.Equal(t, tt.want, sm.s)
		})
	}
}
