package config

import (
	commonconfig "github.com/lyonmu/quebec/app/common/config"
	"github.com/lyonmu/quebec/app/logger"
)

type SystemConfig struct {
	Redis commonconfig.RedisConfig `yaml:"redis" mapstructure:"redis"`
	Auth  AuthConfig               `yaml:"auth" mapstructure:"auth"`
}

type AuthConfig struct {
	Port        uint16                         `yaml:"port" mapstructure:"port"`
	Distributed commonconfig.DistributedConfig `yaml:"distributed" mapstructure:"distributed"`
	Log         logger.LogConfig               `yaml:"log" mapstructure:"log"`
}
