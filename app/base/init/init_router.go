package init

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lyonmu/quebec/app/base/router/v1"
)

func RouterSetup(Router *gin.RouterGroup) {
	v1.RouterV1GroupInstance.UserRouter.InitBasicUserOperationRouterV1(Router)
	v1.RouterV1GroupInstance.SystemRouter.InitSystemMetricsRouterV1(Router)
}
