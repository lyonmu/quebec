package common

type ModuleName string

const (
	ModuleNameGateway ModuleName = "quebec-gateway" // 网关模块-网关服务、服务发现、负载均衡、限流、熔断、认证等
	ModuleNameBasic   ModuleName = "quebec-basic"   // 基础模块-平台中台、基础数据维护等
)
