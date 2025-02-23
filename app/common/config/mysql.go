package config

type MySqlConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`         // 服务器地址
	Port     string `yaml:"port" mapstructure:"port"`         // 端口
	Database string `yaml:"db" mapstructure:"db"`             // 数据库名
	Username string `yaml:"name" mapstructure:"name"`         // 数据库用户名
	Password string `yaml:"password" mapstructure:"password"` // 数据库密码
	TimeOut  uint32 `yaml:"timeout" mapstructure:"timeout"`   // 连接超时时间
}
