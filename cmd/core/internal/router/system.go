package router

import (
	"github.com/gin-gonic/gin"
)

type SystemRouter struct{}

func (r *SystemRouter) InitSystemRouter(Router *gin.RouterGroup) {
	systemRouter := Router.Group("v1/system")
	{
		systemRouter.GET("captcha", systemV1Api.SystemCaptcha)
	}
}
