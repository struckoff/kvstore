package dataitem

type KVDataItem string

func NewKVDataItem(key string) KVDataItem {
	return KVDataItem(key)
}

func (di KVDataItem) ID() string {
	return string(di)
}

func (di KVDataItem) Size() uint64 {
	return 1
}

func (di KVDataItem) Values() []interface{} {
	return []interface{}{string(di)}
}
