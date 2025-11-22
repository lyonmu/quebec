package config

type Jwt struct {
	Sign  string `name:"sign" env:"JWT_SIGN" default:"quebec" help:"jwt签名" mapstructure:"sign" yaml:"sign" json:"sign"`
	Cache int64  `name:"cache" env:"JWT_CACHE" default:"30" help:"验证码缓存时间[小时]" mapstructure:"cache" yaml:"cache" json:"cache"`
}
