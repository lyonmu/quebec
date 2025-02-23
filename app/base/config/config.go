package config

import (
	commonconfig "github.com/lyonmu/quebec/app/common/config"
	"github.com/lyonmu/quebec/app/logger"
)

type SystemConfig struct {
	Redis commonconfig.RedisConfig `yaml:"redis" mapstructure:"redis"`
	Base  BaseConfig               `yaml:"base" mapstructure:"base"`
	Mysql commonconfig.MySqlConfig `yaml:"mysql" mapstructure:"mysql"`
}

type BaseConfig struct {
	Port         uint16                         `yaml:"port" mapstructure:"port"`
	Distributed  commonconfig.DistributedConfig `yaml:"distributed" mapstructure:"distributed"`
	Log          logger.LogConfig               `yaml:"log" mapstructure:"log"`
	RouterPrefix string                         `yaml:"router_prefix" mapstructure:"router_prefix"`
}
