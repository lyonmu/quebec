package common

type ModuleName string

const (
	ModuleNameGateway ModuleName = "quebec-gateway" // 网关模块-网关服务、服务发现、负载均衡、限流、熔断、认证等
	ModuleNameCore    ModuleName = "quebec-core"    // 核心模块-服务发现、负载均衡、限流、熔断、认证等
)
