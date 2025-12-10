package config

import (
	log "github.com/lyonmu/quebec/pkg/logger"
)

type Config struct {
	Version bool          `short:"v" long:"version" help:"版本信息" default:"false" mapstructure:"version" json:"version" yaml:"version"`
	Log     log.LogConfig `embed:"" prefix:"log." mapstructure:"log" json:"log" yaml:"log"`
	Gateway GatewayConfig `embed:"" prefix:"gateway." mapstructure:"gateway" json:"gateway" yaml:"gateway"`
}

type GatewayConfig struct {
	Port   uint16 `name:"port" env:"PORT" default:"59025" help:"端口" mapstructure:"port" yaml:"port" json:"port"`
	Node   int    `name:"node" env:"NODE" default:"1" help:"节点编号" mapstructure:"node" yaml:"node" json:"node"`
	Prefix string `name:"prefix" env:"PREFIX" default:"/gateway/api" help:"路由前缀" mapstructure:"prefix" yaml:"prefix" json:"prefix"`
	Admin  string `name:"admin" env:"ADMIN" default:"127.0.0.1:59024" help:"Admin服务地址" mapstructure:"admin" yaml:"admin" json:"admin"`
}

func (c *Config) MachineID() (int, error) {
	return c.Gateway.Node, nil
}
