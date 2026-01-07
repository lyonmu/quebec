package serializer

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// ProtoCodec Protocol Buffer 编解码器
// 适用于结构化的 proto 消息编解码
type ProtoCodec[K, V any] struct{}

func NewProtoCodec[K, V any]() *ProtoCodec[K, V] {
	return &ProtoCodec[K, V]{}
}

// MarshalKey 实现 Serializer 接口
func (c *ProtoCodec[K, V]) MarshalKey(key *K) ([]byte, error) {
	if key == nil {
		return nil, nil
	}
	// 检查指针类型的 nil
	v := reflect.ValueOf(key)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, nil
	}
	msg, ok := any(key).(proto.Message)
	if !ok {
		return nil, fmt.Errorf("key does not implement proto.Message, got %T", key)
	}
	return proto.Marshal(msg)
}

// MarshalValue 实现 Serializer 接口
func (c *ProtoCodec[K, V]) MarshalValue(value *V) ([]byte, error) {
	if value == nil {
		return nil, nil
	}
	// 检查指针类型的 nil
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, nil
	}
	msg, ok := any(value).(proto.Message)
	if !ok {
		return nil, fmt.Errorf("value does not implement proto.Message, got %T", value)
	}
	return proto.Marshal(msg)
}

// UnmarshalKey 实现 Deserializer 接口
func (c *ProtoCodec[K, V]) UnmarshalKey(raw []byte, key *K) error {
	if key == nil || len(raw) == 0 {
		return nil
	}
	// 检查指针类型的 nil
	v := reflect.ValueOf(key)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	msg, ok := any(key).(proto.Message)
	if !ok {
		return fmt.Errorf("key does not implement proto.Message, got %T", key)
	}
	return proto.Unmarshal(raw, msg)
}

// UnmarshalValue 实现 Deserializer 接口
func (c *ProtoCodec[K, V]) UnmarshalValue(raw []byte, value *V) error {
	if value == nil || len(raw) == 0 {
		return nil
	}
	// 检查指针类型的 nil
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	msg, ok := any(value).(proto.Message)
	if !ok {
		return fmt.Errorf("value does not implement proto.Message, got %T", value)
	}
	return proto.Unmarshal(raw, msg)
}
