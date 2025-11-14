package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type MySQLConfig struct {
	Host     string `name:"host" env:"MYSQL_HOST" default:"127.0.0.1" help:"mysql数据库主机" mapstructure:"host" yaml:"host" json:"host"`
	Port     int    `name:"port" env:"MYSQL_PORT" default:"3306" help:"mysql数据库端口" mapstructure:"port" yaml:"port" json:"port"`
	User     string `name:"user" env:"MYSQL_USER" default:"root" help:"mysql数据库用户" mapstructure:"user" yaml:"user" json:"user"`
	Password string `name:"password" env:"MYSQL_PASSWORD" default:"root" help:"mysql数据库密码" mapstructure:"password" yaml:"password" json:"password"`
	DBName   string `name:"db_name" env:"MYSQL_DB_NAME" default:"quebec" help:"mysql数据库名称" mapstructure:"db_name" yaml:"db_name" json:"db_name"`
}

func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName)
}

type RedisConfig struct {
	Host     []string `name:"host" env:"REDIS_HOST" default:"127.0.0.1:6379" help:"redis服务器地址" mapstructure:"host" yaml:"host" json:"host"`
	Password string   `name:"password" env:"REDIS_PASSWORD" default:"root" help:"redis服务器密码" mapstructure:"password" yaml:"password" json:"password"`
	DB       int      `name:"db" env:"REDIS_DB" default:"1" help:"redis数据库" mapstructure:"db" yaml:"db" json:"db"`
}

func (c *RedisConfig) Client() redis.UniversalClient {
	options := &redis.UniversalOptions{
		Addrs:    c.Host,
		Password: c.Password,
		DB:       c.DB,
		PoolSize: 100,
	}
	return redis.NewUniversalClient(options)
}
