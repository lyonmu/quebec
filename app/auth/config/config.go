package config

import (
	commonconfig "github.com/lyonmu/quebec/app/common/config"
	"github.com/lyonmu/quebec/app/logger"
)

type AuthConfig struct {
	Port  uint16                   `yaml:"port" mapstructure:"port"`
	Log   logger.LogConfig         `yaml:"log" mapstructure:"log"`
	Redis commonconfig.RedisConfig `yaml:"redis" mapstructure:"redis"`
}
