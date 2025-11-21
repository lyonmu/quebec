package constant

type ModuleName string

const (
	ProjectName       ModuleName = "quebec"         // 项目名称
	ModuleNameGateway ModuleName = "quebec-gateway" // 网关模块-xds实现等
	ModuleNameCore    ModuleName = "quebec-core"    // 核心模块-用户管理、数据管理等
)

type YesOrNo int8

const (
	Yes YesOrNo = 1
	No  YesOrNo = 2
)


