package config

import (
	"github.com/lyonmu/quebec/pkg/config"
	log "github.com/lyonmu/quebec/pkg/logger"
)

type Captcha struct {
	Long   int `name:"long" env:"CAPTCHA_LONG" default:"4" help:"验证码长度" mapstructure:"long" yaml:"long" json:"long"`
	Width  int `name:"width" env:"CAPTCHA_WIDTH" default:"240" help:"验证码宽度" mapstructure:"width" yaml:"width" json:"width"`
	Height int `name:"height" env:"CAPTCHA_HEIGHT" default:"80" help:"证码高度" mapstructure:"height" yaml:"height" json:"height"`
}

type CoreConfig struct {
	Captcha Captcha `embed:"" prefix:"captcha." mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Node    int     `name:"node" env:"NODE" default:"1" help:"节点编号" mapstructure:"node" yaml:"node" json:"node"`
}

type Config struct {
	Version bool               `short:"v" long:"version" help:"版本信息" default:"false" mapstructure:"version" json:"version" yaml:"version"`
	Host    string             `long:"host" env:"HOST" help:"服务IP" default:"127.0.0.1" mapstructure:"host" json:"host" yaml:"host"`
	Mode    string             `enum:"local,consul" short:"m" long:"mode" env:"MODE" help:"配置文件模式:local、consul" default:"local" mapstructure:"mode" json:"mode" yaml:"mode"`
	Log     log.LogConfig      `embed:"" prefix:"log." mapstructure:"log" json:"log" yaml:"log"`
	MySQL   config.MySQLConfig `embed:"" prefix:"mysql." mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis   config.RedisConfig `embed:"" prefix:"redis." mapstructure:"redis" json:"redis" yaml:"redis"`
	Core    CoreConfig         `embed:"" prefix:"core." mapstructure:"core" json:"core" yaml:"core"`
}
