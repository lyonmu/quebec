package config

import (
	"github.com/lyonmu/quebec/pkg/config"
	log "github.com/lyonmu/quebec/pkg/logger"
)

type CoreConfig struct {
	Port    uint16  `name:"port" env:"PORT" default:"59024" help:"端口" mapstructure:"port" yaml:"port" json:"port"`
	Node    int     `name:"node" env:"NODE" default:"1" help:"节点编号" mapstructure:"node" yaml:"node" json:"node"`
	Prefix  string  `name:"prefix" env:"PREFIX" default:"/core/api" help:"路由前缀" mapstructure:"prefix" yaml:"prefix" json:"prefix"`
	Captcha Captcha `embed:"" prefix:"captcha." mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Jwt     Jwt     `embed:"" prefix:"jwt." mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}

type Config struct {
	Version bool               `short:"v" long:"version" help:"版本信息" default:"false" mapstructure:"version" json:"version" yaml:"version"`
	Log     log.LogConfig      `embed:"" prefix:"log." mapstructure:"log" json:"log" yaml:"log"`
	MySQL   config.MySQLConfig `embed:"" prefix:"mysql." mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis   config.RedisConfig `embed:"" prefix:"redis." mapstructure:"redis" json:"redis" yaml:"redis"`
	Core    CoreConfig         `embed:"" prefix:"core." mapstructure:"core" json:"core" yaml:"core"`
}

func (c *Config) MachineID() (int, error) {
	return c.Core.Node, nil
}
