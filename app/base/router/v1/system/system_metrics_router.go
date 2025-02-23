package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lyonmu/quebec/app/base/api/v1"
)

type SystemMetricsRouterV1 struct{}

func (*SystemMetricsRouterV1) InitSystemMetricsRouterV1(Router *gin.RouterGroup) {
	systemRouter := Router.Group("/system")
	systemApi := v1.ApiV1GroupInstance.SystemApi
	{
		systemRouter.GET("/metrics", systemApi.Metrics)
	}
}
