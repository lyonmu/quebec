package user

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lyonmu/quebec/app/base/api/v1"
)

type BasicUserOperationRouterV1 struct{}

func (*BasicUserOperationRouterV1) InitBasicUserOperationRouterV1(Router *gin.RouterGroup) {
	baseUserRouterGroup := Router.Group("/user")
	userApi := v1.ApiV1GroupInstance.UserApi
	{
		baseUserRouterGroup.POST("/login", userApi.Login)
		baseUserRouterGroup.POST("/signup", userApi.Signup)
		baseUserRouterGroup.GET("/captcha", userApi.Captcha)
	}
}
