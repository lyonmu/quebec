package mq

import "context"

type Producer[K, V any] interface {
	Produce(ctx context.Context, key *K, payload *V) error
}

type Subscriber[K, V any] interface {
	Subscribe(ctx context.Context, callback func(key *K, payload *V, err error)) error
	Close() error
}

type Decoder[K, V any] interface {
	DecodeKey(raw []byte, key *K) error
	DecodeValue(raw []byte, value *V) error
}

type Encoder[K, V any] interface {
	EncodeKey(key *K) ([]byte, error)
	EncodeValue(value *V) ([]byte, error)
}
