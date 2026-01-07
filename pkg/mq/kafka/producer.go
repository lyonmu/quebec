package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/lyonmu/quebec/pkg/mq"
	ser "github.com/lyonmu/quebec/pkg/mq/serializer"
)

// Producer 泛型生产者实现
type Producer[K, V any] struct {
	producer   sarama.AsyncProducer
	keyCodec   mq.KeySerializer[K]   // Key 编码器
	valueCodec mq.ValueSerializer[V] // Value 编码器
	topic      string
	closeCh    chan struct{}
	wg         sync.WaitGroup
}

// initTopic 是一个独立的辅助函数，用于幂等地创建 Topic
func initTopic(topic string, opts *Options) error {
	admin, err := sarama.NewClusterAdmin(opts.addrs, opts.cfg)
	if err != nil {
		return fmt.Errorf("init topic: create admin failed: %w", err)
	}
	defer admin.Close()

	// 检查 Topic 是否存在
	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("init topic: list topics failed: %w", err)
	}

	if _, exists := topics[topic]; exists {
		return nil
	}

	// 创建 Topic
	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     opts.partitions,
		ReplicationFactor: opts.replicationFactor,
	}, false)

	// 处理 Race Condition
	if err != nil && err != sarama.ErrTopicAlreadyExists {
		return fmt.Errorf("init topic: create topic failed: %w", err)
	}
	return nil
}

// NewProducer 创建一个新的生产者
func NewProducer[K, V any](topic string, opts ...Option) (*Producer[K, V], error) {
	// 1. 初始化配置
	options := defaultOptions()
	for _, apply := range opts {
		apply(options)
	}

	// 2. 核心修复：确保 Topic 存在
	// 如果设置了分区数，则尝试初始化 Topic
	if options.partitions > 0 {
		if err := initTopic(topic, options); err != nil {
			return nil, err
		}
	}

	// 3. 创建 Sarama AsyncProducer
	p, err := sarama.NewAsyncProducer(options.addrs, options.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// 4. 初始化编码器
	var keyCodec mq.KeySerializer[K]
	var valueCodec mq.ValueSerializer[V]

	switch {
	case options.keyCodec != nil && options.valueCodec != nil:
		// 独立配置 key 和 value 编码器
		if kc, ok := options.keyCodec.(mq.KeySerializer[K]); ok {
			keyCodec = kc
		}
		if vc, ok := options.valueCodec.(mq.ValueSerializer[V]); ok {
			valueCodec = vc
		}
	case options.codec != nil:
		// 统一编码器（key 和 value 使用相同编码器）
		if kc, ok := options.codec.(mq.KeySerializer[K]); ok {
			keyCodec = kc
		}
		if vc, ok := options.codec.(mq.ValueSerializer[V]); ok {
			valueCodec = vc
		}
	default:
		// 默认使用 JsonCodec
		keyCodec = ser.NewJsonCodec[K, V]()
		valueCodec = ser.NewJsonCodec[K, V]()
	}

	producer := &Producer[K, V]{
		producer:   p,
		keyCodec:   keyCodec,
		valueCodec: valueCodec,
		topic:      topic,
		closeCh:    make(chan struct{}),
	}

	producer.wg.Add(1)
	go producer.dispatch(options.errorHandler)

	return producer, nil
}

// Produce 实现 Produce 接口
func (p *Producer[K, V]) Produce(ctx context.Context, key *K, payload *V) error {
	// 1. 序列化 Key
	kBytes, err := p.keyCodec.MarshalKey(key)
	if err != nil {
		return fmt.Errorf("marshal key failed: %w", err)
	}

	// 2. 序列化 Value
	vBytes, err := p.valueCodec.MarshalValue(payload)
	if err != nil {
		return fmt.Errorf("marshal value failed: %w", err)
	}

	// 3. 构建消息
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(kBytes),
		Value: sarama.ByteEncoder(vBytes),
	}

	// 4. 发送至 Sarama Input Channel
	select {
	case p.producer.Input() <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.closeCh:
		return fmt.Errorf("producer is closed")
	}
}

// dispatch 负责排空 Error 和 Success 通道，防止死锁
func (p *Producer[K, V]) dispatch(errHandler func(error)) {
	defer p.wg.Done()

	for {
		select {
		case err, ok := <-p.producer.Errors():
			if !ok {
				return
			}
			if errHandler != nil {
				errHandler(err)
			}

		case _, ok := <-p.producer.Successes():
			if !ok {
				return
			}

		case <-p.closeCh:
			return
		}
	}
}

// Close 优雅关闭
func (p *Producer[K, V]) Close() error {
	close(p.closeCh)

	err := p.producer.Close()
	p.wg.Wait()
	return err
}
