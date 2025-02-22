package config

type RedisConfig struct {
	Addr     []string `yaml:"addr" mapstructure:"addr"`
	Password string   `yaml:"password" mapstructure:"password"`
	DB       int      `yaml:"db" mapstructure:"db"`
}

type MongoDBConfig struct {
	Host     string `default:"127.0.0.1" yaml:"host"`  // 主机
	Port     int    `default:"27017" yaml:"port"`      // 端口
	Username string `default:"admin" yaml:"username"`  // 用户名
	Password string `default:"123456" yaml:"password"` // 密码
	AuthDB   string `default:"admin" yaml:"authDB"`    // 认证数据库
}

type DistributedConfig struct {
	ID uint16 `yaml:"id" mapstructure:"id"`
}
