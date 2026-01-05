package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/lyonmu/quebec/pkg/mq"
	"github.com/lyonmu/quebec/pkg/mq/serializer"
)

// Producer 泛型生产者实现
type Producer[K, V any] struct {
	producer sarama.AsyncProducer
	codec    mq.Encode[K, V]
	topic    string
	closeCh  chan struct{}
	wg       sync.WaitGroup
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

	// 4. 初始化 Codec
	codec := serializer.NewJsonCodec[K, V]()

	producer := &Producer[K, V]{
		producer: p,
		codec:    codec,
		topic:    topic,
		closeCh:  make(chan struct{}),
	}

	producer.wg.Add(1)
	go producer.dispatch(options.errorHandler)

	return producer, nil
}

// Produce 实现 Produce 接口
// 这是一个非阻塞操作（除非本地 Buffer 已满）
func (p *Producer[K, V]) Produce(ctx context.Context, key *K, payload *V) error {
	// 1. 序列化 (利用 JsonCodec 的高性能)
	// 注意：EncoderKey/Value 返回的是 []byte
	kBytes, err := p.codec.EncoderKey(key)
	if err != nil {
		return fmt.Errorf("encode key failed: %w", err)
	}

	vBytes, err := p.codec.EncoderValue(payload)
	if err != nil {
		return fmt.Errorf("encode value failed: %w", err)
	}

	// 2. 构建消息
	// sarama.ByteEncoder 只是做类型转换，没有内存拷贝
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(kBytes),
		Value: sarama.ByteEncoder(vBytes),
	}

	// 3. 发送至 Sarama Input Channel
	// 使用 select 监听 context，处理缓冲区满导致的阻塞或超时
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
		// 必须处理 Errors，否则 channel 满了会阻塞 Sender
		case err, ok := <-p.producer.Errors():
			if !ok {
				return
			}
			if errHandler != nil {
				// 将 sarama.ProducerError 解包或直接传递
				errHandler(err)
			}

		// 如果配置了 Return.Successes = true，这里必须读取，否则会死锁
		// 默认配置是 false，这个 case 会被忽略
		case _, ok := <-p.producer.Successes():
			if !ok {
				return
			}
			// 可以在这里做 Metric 统计 (e.g., QPS counter)

		case <-p.closeCh:
			return
		}
	}
}

// Close 优雅关闭
func (p *Producer[K, V]) Close() error {
	close(p.closeCh) // 通知 dispatch 退出

	// 异步关闭 sarama producer，它会负责将 buffer 中的剩余消息刷出
	err := p.producer.Close()

	p.wg.Wait() // 等待 dispatch 处理完剩余的 errors
	return err
}
