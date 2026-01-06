package serializer

import "fmt"

// BinaryCodec 二进制数据透传编解码器
// 性能最优：直接透传二进制数据，最小化内存拷贝
// 适用于 Envoy ALS 等需要直接透传 protobuf 二进制的场景
type BinaryCodec struct{}

func NewBinaryCodec() *BinaryCodec {
	return &BinaryCodec{}
}

// Marshal 编码（实现 Codec 接口）
// 直接返回原始数据的引用，不进行拷贝
// 注意：调用方需确保数据在发送完成前不会被修改
func (c *BinaryCodec) Marshal(key any) ([]byte, error) {
	if key == nil {
		return nil, nil
	}
	b, ok := key.(*[]byte)
	if !ok {
		return nil, fmt.Errorf("BinaryCodec.Marshal requires *[]byte, got %T", key)
	}
	if b == nil || len(*b) == 0 {
		return nil, nil
	}
	return *b, nil
}

// Unmarshal 解码（实现 Codec 接口）
// 直接将 raw 赋值给 key，引用底层数据
func (c *BinaryCodec) Unmarshal(raw []byte, key any) error {
	if key == nil || len(raw) == 0 {
		return nil
	}
	b, ok := key.(*[]byte)
	if !ok {
		return fmt.Errorf("BinaryCodec.Unmarshal requires *[]byte, got %T", key)
	}
	*b = raw
	return nil
}
