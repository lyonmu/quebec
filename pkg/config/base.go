package config

type ConsulConfig struct {
	ConsulUrl       string `name:"url" env:"CONSUL_URL" default:"127.0.0.1:8500" help:"consul 地址" mapstructure:"consul_url" json:"consul_url" yaml:"consul_url"`
	ConsulHttpToken string `name:"token" env:"CONSUL_HTTP_TOKEN" default:"1234567890" help:"consul token" mapstructure:"consul_http_token" json:"consul_http_token" yaml:"consul_http_token"`
	ConsulKey       string `name:"key" env:"CONSUL_KEY" default:"quebec/gateway/config" help:"consul 配置文件key" mapstructure:"consul_key" json:"consul_key" yaml:"consul_key"`
}

type LocalConfig struct {
	ConfigPath string `name:"path" env:"LOCAL_CONFIG_PATH" default:"./config.yaml" help:"本地配置文件路径" mapstructure:"config_path" json:"config_path" yaml:"config_path"`
}
