package serializer

// StringCodec 字符串编解码器
// 简单高效：将字符串直接转换为 UTF-8 字节
// 适用于 Key 为字符串的场景（如用户ID、设备ID）
type StringCodec[K, V string] struct{}

func NewStringCodec[K, V string]() *StringCodec[K, V] {
	return &StringCodec[K, V]{}
}

// MarshalKey 实现 KeySerializer 接口
func (c *StringCodec[K, V]) MarshalKey(key *K) ([]byte, error) {
	if key == nil || *key == "" {
		return nil, nil
	}
	return []byte(*key), nil
}

// MarshalValue 实现 ValueSerializer 接口
func (c *StringCodec[K, V]) MarshalValue(value *V) ([]byte, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	return []byte(*value), nil
}

// UnmarshalKey 实现 Deserializer 接口（用于消费者）
func (c *StringCodec[K, V]) UnmarshalKey(raw []byte, key *K) error {
	if key == nil || len(raw) == 0 {
		return nil
	}
	*key = K(string(raw))
	return nil
}

// UnmarshalValue 实现 Deserializer 接口（用于消费者）
func (c *StringCodec[K, V]) UnmarshalValue(raw []byte, value *V) error {
	if value == nil || len(raw) == 0 {
		return nil
	}
	*value = V(string(raw))
	return nil
}
