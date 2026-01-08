package constant

// ==================== Proxy配置常量 ====================

// 协议类型
type ProxyProtocolType int8

const (
	ProtocolTypeTCP ProxyProtocolType = 1
	ProtocolTypeUDP ProxyProtocolType = 2
)

// 负载均衡策略 (与Envoy一致)
type ProxyLbPolicy int8

const (
	LbPolicyRoundRobin   ProxyLbPolicy = 1 // 轮询 (DEFAULT)
	LbPolicyLeastRequest ProxyLbPolicy = 2 // 最少请求
	LbPolicyRandom       ProxyLbPolicy = 3 // 随机
	LbPolicyRingHash     ProxyLbPolicy = 4 // 环哈希 (一致性哈希)
	LbPolicyMaglev       ProxyLbPolicy = 5 // Maglev (一致性哈希)
)

// HTTP路由匹配类型
type ProxyHttpRouteMatchType int8

const (
	HttpRouteMatchTypePrefix ProxyHttpRouteMatchType = 1 // 前缀匹配
	HttpRouteMatchTypeExact  ProxyHttpRouteMatchType = 2 // 精确匹配
	HttpRouteMatchTypeRegex  ProxyHttpRouteMatchType = 3 // 正则匹配
)

// 证书类型
type CertType int8

const (
	ServerCert CertType = 1 // X509证书（服务器证书+私钥）
	RootCert   CertType = 2 // 根证书（用于验证客户端证书链）
)

// 熔断配置默认值
const (
	DefaultMaxConnections     = 1024  // 最大连接数
	DefaultMaxPendingRequests = 1024  // 最大等待请求数
	DefaultMaxRequests        = 1024  // 最大请求数
	DefaultMaxRetries         = 3     // 最大重试次数
	DefaultConnectTimeoutMs   = 5000  // 连接超时(毫秒)
	DefaultRouteTimeoutMs     = 15000 // 路由超时(毫秒)
	DefaultNumRetries         = 3     // 默认重试次数
	DefaultPerTryTimeoutMs    = 0     // 单次重试超时(毫秒，0表示使用全局超时)
)

// 重试触发条件
type ProxyRetry string

const (
	RetryOn5xx            ProxyRetry = "5xx"
	RetryOnGatewayError   ProxyRetry = "gateway-error"
	RetryOnConnectFailure ProxyRetry = "connect-failure"
	RetryOnRetriable4xx   ProxyRetry = "retriable-4xx"
	RetryOnRefusedStream  ProxyRetry = "refused-stream"
	RetryOnReset          ProxyRetry = "reset"
	RetryOnRetriableCodes ProxyRetry = "retriable-status-codes"
)
