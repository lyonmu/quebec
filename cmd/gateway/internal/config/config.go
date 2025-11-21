package config

import (
	"github.com/lyonmu/quebec/cmd/gateway/internal/common"
	log "github.com/lyonmu/quebec/pkg/logger"
)

type Config struct {
	Version bool          `short:"v" long:"version" help:"版本信息" default:"false" mapstructure:"version" json:"version" yaml:"version"`
	Log     log.LogConfig `embed:"" prefix:"log." mapstructure:"log" json:"log" yaml:"log"`
	Gateway GatewayConfig `embed:"" prefix:"gateway." mapstructure:"gateway" json:"gateway" yaml:"gateway"`
}

type GatewayConfig struct {
	Policy          common.LBPolicy `enum:"1,2,3,4,5" name:"policy" env:"LB_POLICY" default:"1" help:"负载均衡策略 1:随机,2:加权轮询,3:加权最小请求数,4:环形一致性哈希,5:Maglev一致性哈希" mapstructure:"policy" json:"policy" yaml:"policy"`
	Prefix          string          `name:"prefix" env:"ROUTER_PREFIX" default:"/quebec" help:"路由前缀" mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	SvcTag          string          `name:"tag" env:"SVC_TAG" default:"quebec" help:"服务标签"  mapstructure:"tag" json:"tag" yaml:"tag"`
	UpstreamTimeout int             `name:"timeout" env:"UPSTREAM_TIMEOUT" default:"3" help:"上游服务超时时间[秒]" mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	EnableAuth      bool            `name:"auth" env:"ENABLE_AUTH" default:"false" help:"是否启用认证" mapstructure:"auth" json:"auth" yaml:"auth"`
	ProxyPort       uint16          `name:"proxy" env:"PROXY_PORT" default:"11800" help:"代理端口" mapstructure:"proxy_port" json:"proxy_port" yaml:"proxy_port"`
	AdminPort       uint16          `name:"admin" env:"ADMIN_PORT" default:"11801" help:"Envoy Admin服务端口" mapstructure:"admin_port" json:"admin_port" yaml:"admin_port"`
	XdsPort         uint16          `name:"xds" env:"XDS_PORT" default:"11802" help:"XDS服务端口" mapstructure:"xds_port" json:"xds_port" yaml:"xds_port"`
}
