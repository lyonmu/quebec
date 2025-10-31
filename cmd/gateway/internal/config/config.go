package config

import (
	"github.com/lyonmu/quebec/cmd/gateway/internal/common"
	cfg "github.com/lyonmu/quebec/pkg/config"
	log "github.com/lyonmu/quebec/pkg/logger"
)

type Config struct {
	Version  bool             `short:"v" long:"version" help:"版本信息" default:"false" mapstructure:"version" json:"version" yaml:"version"`
	Endpoint string           `short:"e" long:"endpoint" env:"ENDPOINT" help:"服务IP" default:"127.0.0.1:11800" mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	Mode     string           `enum:"local,consul" short:"m" long:"mode" env:"MODE" help:"配置文件模式:local、consul" default:"local" mapstructure:"mode" json:"mode" yaml:"mode"`
	Consul   cfg.ConsulConfig `embed:"" prefix:"consul." mapstructure:"consul" json:"consul" yaml:"consul"`
	Local    cfg.LocalConfig  `embed:"" prefix:"local." mapstructure:"local" json:"local" yaml:"local"`
	Gateway  GatewayConfig    `embed:"" prefix:"gateway." mapstructure:"gateway" json:"gateway" yaml:"gateway"`
	Log      log.LogConfig    `embed:"" prefix:"log." mapstructure:"log" json:"log" yaml:"log"`
}

type GatewayConfig struct {
	LoadBalancerPolicy common.LoadBalancerPolicy `enum:"1,2,3,4,5" name:"policy" env:"LOAD_BALANCER_POLICY" default:"1" help:"负载均衡策略 1:随机,2:加权轮询,3:加权最小请求数,4:环形一致性哈希,5:Maglev一致性哈希"`
}
