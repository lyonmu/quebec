package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lyonmu/quebec/cmd/core/internal/api/http/v1"
	"github.com/lyonmu/quebec/cmd/core/internal/middleware/http"
)

type SystemRouter struct{}

func (r *SystemRouter) InitSystemRouter(Router *gin.RouterGroup, apiGroup v1.V1ApiGroup) {
	systemRouter := Router.Group("v1/system")
	{
		systemRouter.GET("captcha", apiGroup.SystemCaptcha)
		systemRouter.POST("login", apiGroup.SystemLogin)
	}

	systemRouterWithAuth := Router.Group("v1/system", http.JwtAuth())
	{
		systemRouterWithAuth.GET("logout", apiGroup.SystemLogout)
	}
}
