package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lyonmu/quebec/cmd/core/internal/api/http/v1"
)

type SystemRouter struct{}

func (r *SystemRouter) InitSystemRouter(Router *gin.RouterGroup, apiGroup v1.V1ApiGroup) {
	systemRouter := Router.Group("v1/system")
	{
		systemRouter.GET("captcha", apiGroup.SystemCaptcha)
		systemRouter.POST("login", apiGroup.SystemLogin)
	}
}
