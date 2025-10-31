package config

import (
	"github.com/lyonmu/quebec/cmd/gateway/internal/common"
	cfg "github.com/lyonmu/quebec/pkg/config"
)

type Config struct {
	Version          bool   `short:"v" long:"version" help:"版本信息" default:"false"`
	Endpoint         string `short:"e" long:"endpoint" env:"ENDPOINT" help:"服务IP" default:"127.0.0.1:11800"`
	Mode             string `enum:"local,consul" short:"m" long:"mode" env:"MODE" help:"配置文件模式:local、consul" default:"local"`
	cfg.ConsulConfig `embed:"" prefix:"consul."`
	cfg.LocalConfig  `embed:"" prefix:"local."`
	GatewayConfig    `embed:"" prefix:"gateway."`
}

type GatewayConfig struct {
	LoadBalancerPolicy common.LoadBalancerPolicy `enum:"1,2,3,4,5" name:"policy" env:"LOAD_BALANCER_POLICY" default:"1" help:"负载均衡策略 1:随机,2:加权轮询,3:加权最小请求数,4:环形一致性哈希,5:Maglev一致性哈希"`
}
