package config

type ConsulConfig struct {
	ConsulUrl       string `name:"url" env:"CONSUL_URL" default:"127.0.0.1:8500" help:"consul 地址"`
	ConsulHttpToken string `name:"token" env:"CONSUL_HTTP_TOKEN" default:"1234567890" help:"consul token"`
	ConsulKey       string `name:"key" env:"CONSUL_KEY" default:"quebec/gateway/config" help:"consul 配置文件key"`
}

type LocalConfig struct {
	ConfigPath string `name:"path" env:"LOCAL_CONFIG_PATH" default:"./config.yaml" help:"本地配置文件路径"`
}

