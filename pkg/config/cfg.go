package config

type ConsulConfig struct {
	ConsulUrl       string `name:"url" env:"CONSUL_URL" default:"127.0.0.1:8500" help:"consul 地址" mapstructure:"consul_url"`
	ConsulHttpToken string `name:"token" env:"CONSUL_HTTP_TOKEN" default:"1234567890" help:"consul token" mapstructure:"consul_http_token"`
	ConsulKey       string `name:"key" env:"CONSUL_KEY" default:"quebec/gateway/config" help:"consul 配置文件key" mapstructure:"consul_key"`
}

type LocalConfig struct {
	ConfigPath string `name:"path" env:"LOCAL_CONFIG_PATH" default:"./config.yaml" help:"本地配置文件路径" mapstructure:"config_path"`
}

type LogConfig struct {
	LogPath   string `name:"path" env:"LOG_PATH" default:"/var/log/quebec" help:"日志文件路径" mapstructure:"log_path"`                              // 日志文件路径
	LogLevel  string `name:"level" env:"LOG_LEVEL" default:"info" help:"日志级别 [debug, info, warn, error]" mapstructure:"log_level"`             // 日志级别 [debug, info, warn, error]
	LogFormat string `enum:"logfmt,json" name:"format" env:"LOG_FORMAT" default:"logfmt" help:"日志格式 [logfmt, json]" mapstructure:"log_format"` // 日志格式 [logfmt, json]
	LogMax    int64  `name:"max" env:"LOG_MAX" default:"168" help:"日志保留时间(小时)" mapstructure:"log_max"`                                         // 日志保留时间(小时)
}
