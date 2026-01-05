package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

// 默认配置值
const (
	defaultFlushFrequency = 10 * time.Millisecond // 10ms 聚合一次
	defaultFlushBytes     = 16 * 1024             // 16KB 聚合一次
	// defaultChannelSize    = 1024                  // 本地缓冲队列大小
)

type Options struct {
	cfg               *sarama.Config
	addrs             []string
	errorHandler      func(error) // 异步错误回调
	partitions        int32       // 使用 int32 匹配 sarama 定义
	replicationFactor int16       // 副本因子，生产环境建议 >= 2
}

type Option func(*Options)

// defaultOptions 初始化高性能默认值
func defaultOptions() *Options {
	cfg := sarama.NewConfig()

	// --- 高性能核心配置 ---

	// 1. 异步回执：为了极致吞吐，默认不等待成功回调，只处理错误
	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = true

	// 2. 批处理 (Batching)：这是高吞吐的关键
	cfg.Producer.Flush.Frequency = defaultFlushFrequency
	cfg.Producer.Flush.Bytes = defaultFlushBytes
	cfg.Producer.Flush.Messages = 0 // 不限制条数，由 Bytes 或 Frequency 触发

	// 3. 网络优化
	cfg.Net.MaxOpenRequests = 5 // 允许 inflight 请求数 (TCP Pipelining)

	// 4. 可靠性 (根据业务调整，这里取折中方案)
	// cfg.Producer.RequiredAcks = sarama.WaitForLocal // Leader 确认即可，平衡性能与安全
	cfg.Producer.RequiredAcks = sarama.NoResponse // 无需确认，性能最高
	// 若要求数据绝对不丢，请改为 WaitForAll，性能会下降约 30-50%

	// 5. 幂等性 (防止重试导致的数据重复) - 注意：Idempotent 需要 RequiredAcks=WaitForAll
	// 生产环境建议开启以保证 exactly-once 语义，测试环境可禁用以简化配置
	// cfg.Producer.Idempotent = true
	// cfg.Net.MaxOpenRequests = 1 // 开启幂等性时，MaxOpenRequests 必须为 1 (Sarama限制)

	// 6. 分区策略：使用 Hash 策略保证相同 Key 进入同一分区
	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	return &Options{
		cfg:   cfg,
		addrs: []string{"localhost:9092"}, // 默认地址
		errorHandler: func(err error) {
			// 默认打印到 stderr，建议用户覆盖此行为
			// log.Printf("Kafka AsyncProducer Error: %v\n", err)
		},
		partitions:        3, // 默认初始化时使用 Kafka 默认分区策略
		replicationFactor: 2,
	}
}

// WithAddrs 设置 Broker 地址
func WithAddrs(addrs []string) Option {
	return func(o *Options) {
		o.addrs = addrs
	}
}

// WithErrorHandler 设置异步错误处理回调
func WithErrorHandler(h func(error)) Option {
	return func(o *Options) {
		o.errorHandler = h
	}
}

// WithRequiredAcksAll 开启强一致性（牺牲性能）
func WithRequiredAcksAll() Option {
	return func(o *Options) {
		o.cfg.Producer.RequiredAcks = sarama.WaitForAll
		o.cfg.Producer.Idempotent = true // 建议配合幂等性
		o.cfg.Net.MaxOpenRequests = 1
	}
}

func WithSASLGSSAPI(service, krbConf, keytab, relam, username string) Option {
	return func(o *Options) {
		o.cfg.Net.SASL.Enable = true
		o.cfg.Net.SASL.Mechanism = sarama.SASLTypeGSSAPI
		o.cfg.Net.SASL.GSSAPI.ServiceName = service
		o.cfg.Net.SASL.GSSAPI.AuthType = sarama.KRB5_KEYTAB_AUTH
		o.cfg.Net.SASL.GSSAPI.KeyTabPath = keytab
		o.cfg.Net.SASL.GSSAPI.KerberosConfigPath = krbConf
		o.cfg.Net.SASL.GSSAPI.Realm = relam
	}
}

func WithSASLPlaintext(username, password string) Option {
	return func(o *Options) {
		o.cfg.Net.SASL.Enable = true
		o.cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		o.cfg.Net.SASL.User = username
		o.cfg.Net.SASL.Password = password
	}
}

func WithCompressionZSTD() Option {
	return func(o *Options) {
		o.cfg.Producer.Compression = sarama.CompressionZSTD
	}
}

func WithCompressionGZIP() Option {
	return func(o *Options) {
		o.cfg.Producer.Compression = sarama.CompressionGZIP
	}
}

func WithCompressionLZ4() Option {
	return func(o *Options) {
		o.cfg.Producer.Compression = sarama.CompressionLZ4
	}
}

func WithCompressionSnappy() Option {
	return func(o *Options) {
		o.cfg.Producer.Compression = sarama.CompressionSnappy
	}
}

// WithPartitions 设置分区数和副本数
// 如果 replicationFactor <= 0，则默认为 1
func WithPartitions(partitions int32, replicationFactor int16) Option {
	return func(o *Options) {
		o.partitions = partitions
		if replicationFactor > 0 {
			o.replicationFactor = replicationFactor
		} else {
			o.replicationFactor = 1
		}
	}
}

func WithHighThroughput() Option {
	return func(o *Options) {
		o.cfg.Producer.Flush.Frequency = 50 * time.Millisecond
		o.cfg.Producer.Flush.Bytes = 128 * 1024 // 128KB
		o.cfg.Producer.Compression = sarama.CompressionLZ4
	}
}
