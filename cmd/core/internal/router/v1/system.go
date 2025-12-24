package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lyonmu/quebec/cmd/core/internal/api/http/v1"
	"github.com/lyonmu/quebec/cmd/core/internal/middleware/http"
)

type SystemRouter struct{}

func (r *SystemRouter) InitSystemRouter(group *gin.RouterGroup, apiGroup v1.V1ApiGroup) {
	systemRouter := group.Group("v1/system")
	{
		// 验证码
		systemRouter.GET("captcha", apiGroup.SystemCaptcha)
		// 登录
		systemRouter.POST("login", apiGroup.SystemLogin)
	}

	systemRouterWithAuth := group.Group("v1/system", http.JwtAuth())
	{
		// 登出
		systemRouterWithAuth.GET("logout", apiGroup.SystemLogout)
		// 在线用户管理
		systemRouterWithAuth.GET("onlineuser/list", apiGroup.SystemOnlineUserList)
		systemRouterWithAuth.GET("onlineuser/label", apiGroup.SystemOnlineUserLabel)
		systemRouterWithAuth.DELETE("onlineuser/clearance/:id", apiGroup.SystemOnlineUserClearance)
		// 角色管理
		systemRouterWithAuth.GET("role/page", apiGroup.SystemRolePage)
		systemRouterWithAuth.GET("role/list", apiGroup.SystemRoleList)
		systemRouterWithAuth.GET("role/label", apiGroup.SystemRoleLabel)
		systemRouterWithAuth.DELETE("role/:id", apiGroup.SystemRoleDelete)
		systemRouterWithAuth.POST("role", apiGroup.SystemRoleAdd)
		systemRouterWithAuth.PUT("role/:id", apiGroup.SystemRoleEdit)
		systemRouterWithAuth.GET("role/:id", apiGroup.SystemRoleGetById)
		systemRouterWithAuth.PUT("role/enable/:id", apiGroup.SystemRoleEnable)
		// 角色菜单绑定
		systemRouterWithAuth.GET("role/:id/menus", apiGroup.SystemRoleMenus)
		systemRouterWithAuth.PUT("role/:id/menus", apiGroup.SystemRoleBindMenus)
		systemRouterWithAuth.POST("role/:role_id/menu/:menu_id", apiGroup.SystemRoleAddMenu)
		systemRouterWithAuth.DELETE("role/:role_id/menu/:menu_id", apiGroup.SystemRoleRemoveMenu)
		// 菜单管理
		systemRouterWithAuth.GET("menu/page", apiGroup.SystemMenuPage)
		systemRouterWithAuth.GET("menu/list", apiGroup.SystemMenuList)
		systemRouterWithAuth.GET("menu/tree", apiGroup.SystemMenuTree)
		systemRouterWithAuth.GET("menu/label", apiGroup.SystemMenuLabel)
		systemRouterWithAuth.DELETE("menu/:id", apiGroup.SystemMenuDelete)
		systemRouterWithAuth.POST("menu", apiGroup.SystemMenuAdd)
		systemRouterWithAuth.PUT("menu/:id", apiGroup.SystemMenuEdit)
		systemRouterWithAuth.GET("menu/:id", apiGroup.SystemMenuGetById)
		systemRouterWithAuth.PUT("menu/enable/:id", apiGroup.SystemMenuEnable)
		// 用户管理
		systemRouterWithAuth.GET("user/page", apiGroup.SystemUserPage)
		systemRouterWithAuth.GET("user/list", apiGroup.SystemUserList)
		systemRouterWithAuth.GET("user/label", apiGroup.SystemUserLabel)
		systemRouterWithAuth.DELETE("user/:id", apiGroup.SystemUserDelete)
		systemRouterWithAuth.POST("user", apiGroup.SystemUserAdd)
		systemRouterWithAuth.PUT("user/:id", apiGroup.SystemUserEdit)
		systemRouterWithAuth.GET("user/:id", apiGroup.SystemUserGetById)
		systemRouterWithAuth.PUT("user/enable/:id", apiGroup.SystemUserEnable)
		systemRouterWithAuth.PUT("user/password/:id", apiGroup.SystemUserEditPassword)
		systemRouterWithAuth.PUT("user/password/self", apiGroup.SystemUserEditPasswordSelf)
	}
}
