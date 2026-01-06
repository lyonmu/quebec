package mq

import "context"

// Produce 消息生产接口
type Produce[K, V any] interface {
	Produce(ctx context.Context, key *K, payload *V) error
}

// Subscribe 消息订阅接口
type Subscribe[K, V any] interface {
	Subscribe(ctx context.Context, callback func(key *K, payload *V, err error)) error
	Close() error
}

// Serializer 编码接口（生产者发送时使用）
type Serializer[K, V any] interface {
	MarshalKey(key *K) ([]byte, error)
	MarshalValue(value *V) ([]byte, error)
}

// Deserializer 解码接口（消费者读取时使用）
type Deserializer[K, V any] interface {
	UnmarshalKey(raw []byte, key *K) error
	UnmarshalValue(raw []byte, value *V) error
}

// Codec 通用编解码器接口（非泛型，支持 ProtoCodec、BinaryCodec）
type Codec interface {
	Marshal(key any) ([]byte, error)
	Unmarshal(raw []byte, key any) error
}

// 兼容别名（可选，逐步迁移）
// Deprecated: 使用 Serializer 替代
type Encode[K, V any] = Serializer[K, V]

// Deprecated: 使用 Deserializer 替代
type Decode[K, V any] = Deserializer[K, V]
