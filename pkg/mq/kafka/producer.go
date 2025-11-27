package kafka

import (
	"context"
	"github.com/lyonmu/quebec/pkg/mq"
	"github.com/lyonmu/quebec/pkg/mq/codec"
	"time"

	"github.com/IBM/sarama"
)

type producer[K, V any] struct {
	topic         string
	asyncWriter   sarama.AsyncProducer
	kvPairEncoder mq.Encoder[K, V]
}

func NewJSONProducer[K, V any](
	addrs []string,
	topic string,
	opts ...Option,
) (mq.Producer[K, V], error) {
	jsonEncoder := codec.NewJsonEncoder[K, V]()
	return NewProducer(addrs, topic, jsonEncoder, opts...)
}

func NewProducer[K, V any](
	addrs []string,
	topic string,
	kvPairEncoder mq.Encoder[K, V],
	opts ...Option,
) (mq.Producer[K, V], error) {
	cfg := sarama.NewConfig()
	cfg.Producer.MaxMessageBytes = 20 * 1024 * 1024
	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = false
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Flush.Messages = 1000
	cfg.Producer.Flush.Frequency = 1 * time.Second
	cfg.Producer.Flush.MaxMessages = 10000
	cfg.ChannelBufferSize = 500
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	if kvPairEncoder == nil {
		return nil, mq.ErrNilEncoder
	}

	o := &Options{
		cfg: cfg,
	}

	for _, opt := range opts {
		opt(o)
	}

	writer, err := sarama.NewAsyncProducer(addrs, o.cfg)
	if err != nil {
		return nil, err
	}

	return &producer[K, V]{
		topic:         topic,
		asyncWriter:   writer,
		kvPairEncoder: kvPairEncoder,
	}, nil
}

func NewRawProducer(
	addrs []string,
	topic string,
	opts ...Option,
) (mq.Producer[[]byte, []byte], error) {
	cfg := sarama.NewConfig()
	cfg.Producer.MaxMessageBytes = 20 * 1024 * 1024
	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = false
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Flush.Messages = 1000
	cfg.Producer.Flush.Frequency = 1 * time.Second
	cfg.Producer.Flush.MaxMessages = 10000
	cfg.ChannelBufferSize = 500
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	o := &Options{
		cfg: cfg,
	}

	for _, opt := range opts {
		opt(o)
	}

	writer, err := sarama.NewAsyncProducer(addrs, o.cfg)
	if err != nil {
		return nil, err
	}

	return &producer[[]byte, []byte]{
		topic:         topic,
		asyncWriter:   writer,
		kvPairEncoder: codec.NewRawEncoder(),
	}, nil
}

func (p *producer[K, V]) Produce(ctx context.Context, key *K, payload *V) error {
	if payload == nil {
		return nil
	}
	value, err := p.kvPairEncoder.EncodeValue(payload)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(value),
	}
	if key != nil {
		k, err := p.kvPairEncoder.EncodeKey(key)
		if err != nil {
			return err
		}
		msg.Key = sarama.ByteEncoder(k)
	}
	input := p.asyncWriter.Input()
	input <- msg
	return nil
}
