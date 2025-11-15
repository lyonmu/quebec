package config

import (
	"database/sql"
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

// DSN 返回mysql数据库连接字符串
func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName)
}

// EnsureDatabase 检查数据库是否存在，如果不存在则创建
func (c *MySQLConfig) EnsureDatabase() error {
	// 连接到 MySQL 服务器（不指定数据库名）
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local",
		c.User, c.Password, c.Host, c.Port)

	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to connect to mysql server: %v", err)
	}
	defer db.Close()

	// 检查数据库是否存在
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?)"
	err = db.QueryRow(query, c.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %v", err)
	}

	// 如果数据库不存在，则创建
	if !exists {
		createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", c.DBName)
		_, err = db.Exec(createSQL)
		if err != nil {
			return fmt.Errorf("failed to create database %s: %v", c.DBName, err)
		}
	}

	return nil
}

type RedisConfig struct {
	Host     []string `name:"host" env:"REDIS_HOST" default:"127.0.0.1:6379" help:"redis服务器地址" mapstructure:"host" yaml:"host" json:"host"`
	Password string   `name:"password" env:"REDIS_PASSWORD" default:"root" help:"redis服务器密码" mapstructure:"password" yaml:"password" json:"password"`
	DB       int      `name:"db" env:"REDIS_DB" default:"1" help:"redis数据库" mapstructure:"db" yaml:"db" json:"db"`
}

// Client 返回redis客户端
func (c *RedisConfig) Client(name string) redis.UniversalClient {
	options := &redis.UniversalOptions{
		Addrs:                 c.Host,
		Password:              c.Password,
		DB:                    c.DB,
		ClientName:            name,
		ContextTimeoutEnabled: true,
		PoolFIFO:              false,
	}
	return redis.NewUniversalClient(options)
}
