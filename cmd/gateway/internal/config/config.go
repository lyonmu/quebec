package config

import cfg "github.com/lyonmu/quebec/pkg/config"

type Config struct {
	Version          bool   `short:"v" long:"version" help:"版本信息" default:"false"`
	Endpoint         string `short:"e" long:"endpoint" env:"ENDPOINT" help:"服务IP" default:"127.0.0.1:11800"`
	Mode             string `enum:"local,consul" short:"m" long:"mode" env:"MODE" help:"配置文件模式:local、consul" default:"local"`
	cfg.ConsulConfig `embed:"" prefix:"consul."`
	cfg.LocalConfig  `embed:"" prefix:"local."`
}
