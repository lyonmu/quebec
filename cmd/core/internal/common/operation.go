package common

type OperationType int

const (
	OperationLogin               OperationType = 1  // 登录操作
	OperationLogout              OperationType = 2  // 登出操作
	OperationUserCreate          OperationType = 3  // 创建用户
	OperationUserUpdate          OperationType = 4  // 更新用户
	OperationUserDelete          OperationType = 5  // 删除用户
	OperationRoleCreate          OperationType = 6  // 创建角色
	OperationRoleUpdate          OperationType = 7  // 更新角色
	OperationRoleDelete          OperationType = 8  // 删除角色
	OperationMenuCreate          OperationType = 9  // 创建菜单
	OperationMenuUpdate          OperationType = 10 // 更新菜单
	OperationMenuDelete          OperationType = 11 // 删除菜单
	OperationOnlineUserClearance OperationType = 12 // 踢出在线用户
	OperationLogQuery            OperationType = 13 // 查询操作日志
	OperationLogExport           OperationType = 14 // 导出操作日志
	OperationLogClear            OperationType = 15 // 清理操作日志
	OperationPasswordChange      OperationType = 16 // 修改密码
	OperationUserEnable          OperationType = 17 // 启用/禁用用户
	OperationRoleEnable          OperationType = 18 // 启用/禁用角色
	OperationMenuEnable          OperationType = 19 // 启用/禁用菜单
	OperationRoleBindMenus       OperationType = 20 // 角色绑定菜单
)
