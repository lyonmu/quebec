package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lyonmu/quebec/cmd/core/internal/api/http/v1"
	"github.com/lyonmu/quebec/cmd/core/internal/api/http/v1/system"
	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/middleware/http"
	middlewarehttp "github.com/lyonmu/quebec/cmd/core/internal/middleware/http"
)

type SystemRouter struct{}

func (r *SystemRouter) InitSystemRouter(group *gin.RouterGroup, apiGroup v1.V1ApiGroup) {
	// 创建操作日志中间件
	operationLogMiddleware := middlewarehttp.NewOperationLogMiddleware(system.GetSystemSvc())

	systemRouter := group.Group("v1/system")
	{
		// 验证码
		systemRouter.GET("captcha", apiGroup.SystemCaptcha)
		// 登录（需要记录操作日志）
		systemRouter.POST("login", operationLogMiddleware.Handle(common.OperationLogin), apiGroup.SystemLogin)
	}

	systemRouterWithAuth := group.Group("v1/system", http.JwtAuth())
	{
		// 登出（需要记录操作日志）
		systemRouterWithAuth.GET("logout", operationLogMiddleware.Handle(common.OperationLogout), apiGroup.SystemLogout)

		// === 在线用户管理 ===
		systemRouterWithAuth.GET("onlineuser/list", apiGroup.SystemOnlineUserList)
		systemRouterWithAuth.GET("onlineuser/label", apiGroup.SystemOnlineUserLabel)
		// 踢出在线用户（需要记录操作日志）
		systemRouterWithAuth.DELETE("onlineuser/clearance/:id", operationLogMiddleware.Handle(common.OperationOnlineUserClearance), apiGroup.SystemOnlineUserClearance)

		// === 角色管理 ===
		systemRouterWithAuth.GET("role/page", apiGroup.SystemRolePage)
		systemRouterWithAuth.GET("role/list", apiGroup.SystemRoleList)
		systemRouterWithAuth.GET("role/label", apiGroup.SystemRoleLabel)
		// 创建角色（需要记录操作日志）
		systemRouterWithAuth.POST("role", operationLogMiddleware.Handle(common.OperationRoleCreate), apiGroup.SystemRoleAdd)
		// 更新角色（需要记录操作日志）
		systemRouterWithAuth.PUT("role/:id", operationLogMiddleware.Handle(common.OperationRoleUpdate), apiGroup.SystemRoleEdit)
		// 删除角色（需要记录操作日志）
		systemRouterWithAuth.DELETE("role/:id", operationLogMiddleware.Handle(common.OperationRoleDelete), apiGroup.SystemRoleDelete)
		systemRouterWithAuth.GET("role/:id", apiGroup.SystemRoleGetById)
		// 启用/禁用角色（需要记录操作日志）
		systemRouterWithAuth.PUT("role/enable/:id", operationLogMiddleware.Handle(common.OperationRoleEnable), apiGroup.SystemRoleEnable)
		// 角色菜单绑定（需要记录操作日志）
		systemRouterWithAuth.PUT("role/:id/menus", operationLogMiddleware.Handle(common.OperationRoleBindMenus), apiGroup.SystemRoleBindMenus)

		// 角色菜单绑定相关接口（不需要记录日志 - 查询操作）
		systemRouterWithAuth.GET("role/:id/menus", apiGroup.SystemRoleMenus)
		systemRouterWithAuth.POST("role/:id/menu/:menu_id", apiGroup.SystemRoleAddMenu)
		systemRouterWithAuth.DELETE("role/:id/menu/:menu_id", apiGroup.SystemRoleRemoveMenu)

		// === 菜单管理 ===
		systemRouterWithAuth.GET("menu/page", apiGroup.SystemMenuPage)
		systemRouterWithAuth.GET("menu/list", apiGroup.SystemMenuList)
		systemRouterWithAuth.GET("menu/tree", apiGroup.SystemMenuTree)
		systemRouterWithAuth.GET("menu/label", apiGroup.SystemMenuLabel)
		systemRouterWithAuth.GET("menu/:id", apiGroup.SystemMenuGetById)
		// 启用/禁用菜单（需要记录操作日志）
		systemRouterWithAuth.PUT("menu/enable/:id", operationLogMiddleware.Handle(common.OperationMenuEnable), apiGroup.SystemMenuEnable)

		// === 用户管理 ===
		systemRouterWithAuth.GET("user/page", apiGroup.SystemUserPage)
		systemRouterWithAuth.GET("user/list", apiGroup.SystemUserList)
		systemRouterWithAuth.GET("user/label", apiGroup.SystemUserLabel)
		// 创建用户（需要记录操作日志）
		systemRouterWithAuth.POST("user", operationLogMiddleware.Handle(common.OperationUserCreate), apiGroup.SystemUserAdd)
		// 更新用户（需要记录操作日志）
		systemRouterWithAuth.PUT("user/:id", operationLogMiddleware.Handle(common.OperationUserUpdate), apiGroup.SystemUserEdit)
		// 删除用户（需要记录操作日志）
		systemRouterWithAuth.DELETE("user/:id", operationLogMiddleware.Handle(common.OperationUserDelete), apiGroup.SystemUserDelete)
		systemRouterWithAuth.GET("user/:id", apiGroup.SystemUserGetById)
		// 启用/禁用用户（需要记录操作日志）
		systemRouterWithAuth.PUT("user/enable/:id", operationLogMiddleware.Handle(common.OperationUserEnable), apiGroup.SystemUserEnable)
		// 修改密码（需要记录操作日志）
		systemRouterWithAuth.PUT("user/password/:id", operationLogMiddleware.Handle(common.OperationPasswordChange), apiGroup.SystemUserEditPassword)
		systemRouterWithAuth.PUT("user/password/self", operationLogMiddleware.Handle(common.OperationPasswordChange), apiGroup.SystemUserEditPasswordSelf)

		// === 操作日志 ===
		systemRouterWithAuth.GET("operation-log/page", apiGroup.SystemOperationLogPage)
	}
}
