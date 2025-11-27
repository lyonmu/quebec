package kafka

import (
	"bytes"
	"context"
	"github.com/lyonmu/quebec/pkg/mq"
	"github.com/lyonmu/quebec/pkg/mq/codec"
	"reflect"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type subscriber[K, V any] struct {
	client sarama.Client
	topic  string
	group  string

	reader *reader[K, V]
}

func NewJSONSubscriberWithKey[K, V any](addrs []string, topic, group string, reader *reader[K, V], opts ...Option) (mq.Subscriber[K, V], error) {
	jsonDecoder := codec.NewJsonDecoderWithKey[K, V]()
	reader.decoder = jsonDecoder
	return NewSubscriber(addrs, topic, group, reader, opts...)
}

func NewJSONSubscriber[K, V any](addrs []string, topic, group string, reader *reader[K, V], opts ...Option) (mq.Subscriber[K, V], error) {
	jsonDecoder := codec.NewJsonDecoder[K, V]()
	reader.decoder = jsonDecoder
	return NewSubscriber(addrs, topic, group, reader, opts...)
}

func NewSubscriber[K, V any](addrs []string, topic, group string, reader *reader[K, V], opts ...Option) (mq.Subscriber[K, V], error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = false
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = 300 * time.Millisecond
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Group.Session.Timeout = 10 * time.Second
	cfg.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	cfg.Consumer.Fetch.Default = 1 * 1024 * 1024
	cfg.Consumer.MaxWaitTime = 150 * time.Millisecond
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	o := &Options{
		cfg:        cfg,
		partitions: 1,
	}

	for _, opt := range opts {
		opt(o)
	}
	client, err := sarama.NewClient(addrs, o.cfg)
	if err != nil {
		return nil, err
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return nil, err
	}

	topics, err := admin.ListTopics()
	if err != nil {
		return nil, err
	}

	detail := &sarama.TopicDetail{
		NumPartitions:     int32(o.partitions),
		ReplicationFactor: 1,
	}
	if _, exists := topics[topic]; !exists {
		err := admin.CreateTopic(topic, detail, false)
		if err != nil {
			return nil, err
		}
	} else {
		metadata, err := admin.DescribeTopics([]string{topic})
		if err != nil {
			return nil, err
		}
		if len(metadata[0].Partitions) != int(o.partitions) {
			err = admin.DeleteTopic(topic)
			if err != nil {
				return nil, err
			}
			err = admin.CreateTopic(topic, detail, false)
			if err != nil {
				return nil, err
			}
		}
	}

	return &subscriber[K, V]{
		client: client,
		topic:  topic,
		group:  group,
		reader: reader,
	}, nil
}

func (s *subscriber[K, V]) Subscribe(ctx context.Context, callback func(key *K, payload *V, err error)) error {
	if s.reader.decoder == nil {
		return mq.ErrNilDecoder
	}
	s.reader.callback = callback

	group, err := sarama.NewConsumerGroupFromClient(s.group, s.client)
	if err != nil {
		return err
	}

	for {
		if err := group.Consume(ctx, []string{s.topic}, s.reader); err != nil {
			continue
		}
	}
}

func (s *subscriber[K, V]) Close() error {
	return s.client.Close()
}

type reader[K, V any] struct {
	keyPool   *sync.Pool
	valuePool *sync.Pool

	newKV func() (*K, *V)

	decoder mq.Decoder[K, V]

	callback func(key *K, payload *V, err error)
}

func NewReader[K, V any](newKV func() (*K, *V), decoder mq.Decoder[K, V]) *reader[K, V] {
	return &reader[K, V]{
		newKV:   newKV,
		decoder: decoder,
	}
}

func NewPoolReader[K, V any](decoder ...mq.Decoder[K, V]) *reader[K, V] {
	reader := &reader[K, V]{
		keyPool: &sync.Pool{
			New: func() any {
				return new(K)
			},
		},
		valuePool: &sync.Pool{
			New: func() any {
				return new(V)
			},
		},
	}
	reader.newKV = func() (*K, *V) {
		return reader.keyPool.Get().(*K), reader.valuePool.Get().(*V)
	}
	if len(decoder) > 0 {
		reader.decoder = decoder[0]
	}
	return reader
}

type Bytes []byte

func NewRawReader() *reader[bytes.Buffer, bytes.Buffer] {
	return NewPoolReader(codec.NewRawDecoder())
}

func (r *reader[K, V]) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (r *reader[K, V]) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (r *reader[K, V]) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")
		k, v := r.newKV()

		err := r.decoder.DecodeKey(msg.Key, k)
		if err != nil {
			r.callback(nil, nil, err)
			r.reset(k, v)
			continue
		}
		err = r.decoder.DecodeValue(msg.Value, v)
		if err != nil {
			r.callback(nil, nil, err)
			r.reset(k, v)
			continue
		}

		r.callback(k, v, nil)
		r.reset(k, v)
	}

	return nil
}

type object interface {
	Reset()
}

func (r *reader[K, V]) reset(k *K, v *V) {
	if k != nil {
		if k, ok := reflect.ValueOf(k).Interface().(object); ok {
			k.Reset()
			if r.keyPool != nil {
				r.keyPool.Put(k)
			}
		}
	}
	if v != nil {
		if v, ok := reflect.ValueOf(v).Interface().(object); ok {
			v.Reset()
			if r.valuePool != nil {
				r.valuePool.Put(v)
			}
		}
	}
}
