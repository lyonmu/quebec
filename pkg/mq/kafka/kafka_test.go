package kafka

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// 定义集成测试专用的 payload
type IntegrationPayload struct {
	MsgID     int    `json:"msg_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"ts"`
}

func TestProducer_RealCluster_Throughput(t *testing.T) {
	// 1. 定义集群地址 (根据你提供的真实信息)
	brokers := []string{
		"kafka-svc-1:9192",
		"kafka-svc-2:9292",
		"kafka-svc-3:9392",
	}

	topic := "integration-test-topic-v1"
	totalMessages := 1000 // 测试发送 1000 条

	// 用于统计异步错误的计数器
	var errCount int64

	// 2. 初始化 Producer
	// 架构配置：3 节点集群 -> 3 副本高可用，6 分区负载均衡
	p, err := NewProducer[string, IntegrationPayload](topic,
		WithAddrs(brokers),
		WithPartitions(6, 3), // 自动创建 Topic: 6分区/3副本
		WithCompressionLZ4(), // 开启 LZ4 压缩
		WithErrorHandler(func(e error) {
			atomic.AddInt64(&errCount, 1)
			t.Logf("[Async Error] %v", e)
		}),
		WithSASLPlaintext("mu", "lyonmu"),
	)

	if err != nil {
		t.Fatalf("Failed to connect to cluster %v: %v", brokers, err)
	}
	defer func() {
		// Close 会强制 flush buffer 中的剩余消息
		if err := p.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	t.Logf("Successfully connected to Kafka cluster: %v", brokers)

	// 3. 并发发送测试
	ctx := context.Background()
	start := time.Now()

	for i := 0; i < totalMessages; i++ {
		key := fmt.Sprintf("k-%d", i)
		payload := IntegrationPayload{
			MsgID:     i,
			Content:   "Hello Kafka Integration Test",
			Timestamp: time.Now().UnixNano(),
		}

		// Produce 是异步入队，速度极快
		if err := p.Produce(ctx, &key, &payload); err != nil {
			t.Fatalf("Enqueue failed at index %d: %v", i, err)
		}
	}

	// 4. 观察与验证
	duration := time.Since(start)
	t.Logf("Enqueued %d messages in %v (TPS: %.2f)",
		totalMessages, duration, float64(totalMessages)/duration.Seconds())

	// 给一点时间让后台 Goroutine 处理网络 IO 和回调
	// 注意：在真实业务中 p.Close() 会自动处理等待，这里 sleep 只是为了在 Close 前有机会看到 log
	time.Sleep(2 * time.Second)

	// 5. 验证结果
	finalErrs := atomic.LoadInt64(&errCount)
	if finalErrs > 0 {
		t.Errorf("Test completed with %d async errors (see logs above)", finalErrs)
	} else {
		t.Logf("Test Passed: 0 async errors detected.")
	}
}
