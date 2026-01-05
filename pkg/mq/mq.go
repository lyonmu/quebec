package mq

import "context"

type Produce[K, V any] interface {
	Produce(ctx context.Context, key *K, payload *V) error
}

type Subscribe[K, V any] interface {
	Subscribe(ctx context.Context, callback func(key *K, payload *V, err error)) error
	Close() error
}

type Encode[K, V any] interface {
	EncoderKey(key *K) ([]byte, error)
	EncoderValue(value *V) ([]byte, error)
}
type Decode[K, V any] interface {
	DecoderKey(raw []byte, key *K) error
	DecoderValue(raw []byte, value *V) error
}


