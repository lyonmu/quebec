package serializer

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// ProtoCodec Protocol Buffer 编解码器
// 适用于结构化的 proto 消息编解码
type ProtoCodec struct{}

func NewProtoCodec() *ProtoCodec {
	return &ProtoCodec{}
}

// Marshal 编码（实现 Codec 接口）
// key 必须是 proto.Message 类型
func (c *ProtoCodec) Marshal(key any) ([]byte, error) {
	if key == nil {
		return nil, nil
	}
	// 检查指针类型的 nil
	v := reflect.ValueOf(key)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, nil
	}
	msg, ok := key.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("key does not implement proto.Message, got %T", key)
	}
	return proto.Marshal(msg)
}

// Unmarshal 解码（实现 Codec 接口）
// key 必须是 proto.Message 类型的指针
func (c *ProtoCodec) Unmarshal(raw []byte, key any) error {
	if key == nil || len(raw) == 0 {
		return nil
	}
	// 检查指针类型的 nil
	v := reflect.ValueOf(key)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	msg, ok := key.(proto.Message)
	if !ok {
		return fmt.Errorf("key does not implement proto.Message, got %T", key)
	}
	return proto.Unmarshal(raw, msg)
}
