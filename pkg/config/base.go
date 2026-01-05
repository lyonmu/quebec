package config

type ConsulConfig struct {
	ConsulUrl       string `name:"url" env:"CONSUL_URL" default:"127.0.0.1:8500" help:"consul 地址" mapstructure:"consul_url" json:"consul_url" yaml:"consul_url"`
	ConsulHttpToken string `name:"token" env:"CONSUL_HTTP_TOKEN" default:"1234567890" help:"consul token" mapstructure:"consul_http_token" json:"consul_http_token" yaml:"consul_http_token"`
	ConsulKey       string `name:"key" env:"CONSUL_KEY" default:"quebec/gateway/config" help:"consul 配置文件key" mapstructure:"consul_key" json:"consul_key" yaml:"consul_key"`
}

type KafkaSASLConfig struct {
	Enable    bool   `name:"enable" env:"KAFKA_SASL_ENABLE" default:"false" help:"是否启用 kafka sasl 验证" mapstructure:"enable" json:"enable" yaml:"enable"`
	Mechanism string `enum:"PLAIN,SCRAM-SHA-256,SCRAM-SHA-512" name:"mechanism" env:"KAFKA_SASL_MECHANISM" default:"PLAIN" help:"kafka sasl 认证机制 [可选PLAIN,SCRAM-SHA-256,SCRAM-SHA-512]" mapstructure:"mechanism" json:"mechanism" yaml:"mechanism"`
	Username  string `name:"username" env:"KAFKA_SASL_USERNAME" default:"root" help:"kafka 用户名" mapstructure:"username" json:"username" yaml:"username"`
	Password  string `name:"password" env:"KAFKA_SASL_PASSWORD" default:"root" help:"kafka 密码" mapstructure:"password" json:"password" yaml:"password"`
}
type KafkaConfig struct {
	Brokers     []string        `name:"brokers" env:"KAFKA_BROKERS" default:"localhost:9092" help:"kafka broker 列表，逗号分隔" mapstructure:"brokers" json:"brokers" yaml:"brokers"`
	SASL        KafkaSASLConfig `embed:"" prefix:"sasl." mapstructure:"sasl" json:"sasl" yaml:"sasl"`
	Partitions  int             `name:"partitions" env:"KAFKA_PARTITIONS" default:"3" help:"默认分区数" mapstructure:"partitions" json:"partitions" yaml:"partitions"`
	Replication int             `name:"replication" env:"KAFKA_REPLICATION" default:"2" help:"默认副本数" mapstructure:"replication" json:"replication" yaml:"replication"`
	AlsTopic    string          `name:"als_topic" env:"KAFKA_ALS_TOPIC" default:"envoy_als_events" help:"envoy access log 主题名称" mapstructure:"als_topic" json:"als_topic" yaml:"als_topic"`
}
