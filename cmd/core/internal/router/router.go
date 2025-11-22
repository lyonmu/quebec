package router

import (
	"github.com/gin-gonic/gin"
	v1Api "github.com/lyonmu/quebec/cmd/core/internal/api/http/v1"
	_ "github.com/lyonmu/quebec/cmd/core/internal/docs"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	v1Route "github.com/lyonmu/quebec/cmd/core/internal/router/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	v1api   = v1Api.V1ApiGroup{}
	v1route = v1Route.V1Router{}
)

type RouterGroup struct {
	v1Api.V1ApiGroup
}

func InitRouter(e *gin.Engine) {

	// Router group
	routerGroup := e.Group(global.Cfg.Core.Prefix)

	routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Init system router
	v1route.InitSystemRouter(routerGroup, v1api)

	global.Logger.Sugar().Info("router register success")

}
