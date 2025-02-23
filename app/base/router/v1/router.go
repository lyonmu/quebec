package v1

import (
	"github.com/lyonmu/quebec/app/base/router/v1/system"
	"github.com/lyonmu/quebec/app/base/router/v1/user"
)

type RouterV1Group struct {
	UserRouter   user.BasicUserOperationRouterV1
	SystemRouter system.SystemMetricsRouterV1
}

var RouterV1GroupInstance = new(RouterV1Group)
