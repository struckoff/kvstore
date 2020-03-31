package kvstore

// DataItem represents string key as balancer item
type DataItem string

func (di DataItem) ID() string   { return string(di) }
func (di DataItem) Size() uint64 { return 1 }
func (di DataItem) Values() []interface{} {
	vals := make([]interface{}, 1)
	vals[0] = string(di)
	return vals
}
