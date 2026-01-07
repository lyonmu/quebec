package serializer

import "fmt"

// BinaryCodec 二进制数据透传编解码器
// 性能最优：直接透传二进制数据，最小化内存拷贝
// 适用于 Envoy ALS 等需要直接透传 protobuf 二进制的场景
type BinaryCodec[K, V any] struct{}

func NewBinaryCodec[K, V any]() *BinaryCodec[K, V] {
	return &BinaryCodec[K, V]{}
}

// MarshalKey 实现 Serializer 接口
func (c *BinaryCodec[K, V]) MarshalKey(key *K) ([]byte, error) {
	if key == nil {
		return nil, nil
	}
	b, ok := any(key).(*[]byte)
	if !ok {
		return nil, fmt.Errorf("BinaryCodec.MarshalKey requires *[]byte, got %T", key)
	}
	if b == nil || len(*b) == 0 {
		return nil, nil
	}
	return *b, nil
}

// MarshalValue 实现 Serializer 接口
func (c *BinaryCodec[K, V]) MarshalValue(value *V) ([]byte, error) {
	if value == nil {
		return nil, nil
	}
	b, ok := any(value).(*[]byte)
	if !ok {
		return nil, fmt.Errorf("BinaryCodec.MarshalValue requires *[]byte, got %T", value)
	}
	if b == nil || len(*b) == 0 {
		return nil, nil
	}
	return *b, nil
}

// UnmarshalKey 实现 Deserializer 接口
func (c *BinaryCodec[K, V]) UnmarshalKey(raw []byte, key *K) error {
	if key == nil || len(raw) == 0 {
		return nil
	}
	b, ok := any(key).(*[]byte)
	if !ok {
		return fmt.Errorf("BinaryCodec.UnmarshalKey requires *[]byte, got %T", key)
	}
	*b = raw
	return nil
}

// UnmarshalValue 实现 Deserializer 接口
func (c *BinaryCodec[K, V]) UnmarshalValue(raw []byte, value *V) error {
	if value == nil || len(raw) == 0 {
		return nil
	}
	b, ok := any(value).(*[]byte)
	if !ok {
		return fmt.Errorf("BinaryCodec.UnmarshalValue requires *[]byte, got %T", value)
	}
	*b = raw
	return nil
}
