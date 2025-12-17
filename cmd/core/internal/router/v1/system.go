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
	}
}
