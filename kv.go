package page

type KeyValue struct {
	Key   string
	Value string
}

func NewKeyValue(key string, val string) *KeyValue {
	return &KeyValue{
		Key:   key,
		Value: val,
	}
}

func (kv *KeyValue) GetKeySize() uint32 {
	return len(kv.Key)
}

func (kv *KeyValue) GetValueSize() uint32 {
	return len(kv.Value)
}

func (kv *KeyValue) GetSize() uint32 {
	return kv.GetKeySize() + kv.GetValueSize()
}
