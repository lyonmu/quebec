package common

type ModuleName string

const (
	ProjectName       ModuleName = "quebec"         // 项目名称
	ModuleNameGateway ModuleName = "quebec-gateway" // 网关模块-网关服务、服务发现、负载均衡、限流、熔断、认证等
	ModuleNameCore    ModuleName = "quebec-core"    // 核心模块-服务发现、负载均衡、限流、熔断、认证等
)

type YesOrNo int8

const (
	Yes YesOrNo = 1
	No  YesOrNo = 2
)

type MenuType int8

const (
	MenuTypeDirectory MenuType = 1
	MenuTypeMenu      MenuType = 2
	MenuTypeButton    MenuType = 3
)

type DataRelationshipType int8

const (
	// 角色与菜单
	DataRelationshipTypeRoleToMenu DataRelationshipType = 1
	// 用户与角色
	DataRelationshipTypeUserToRole DataRelationshipType = 2
	// 菜单与按钮
	DataRelationshipTypeManyToMany DataRelationshipType = 3
)
