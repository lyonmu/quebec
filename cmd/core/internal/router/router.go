package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lyonmu/quebec/cmd/core/internal/docs"
	"github.com/lyonmu/quebec/cmd/core/internal/api/http/v1/system"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	systemV1Api = system.SystemV1ApiGroup{}
)

type RouterGroup struct {
	SystemRouter
}

var router = new(RouterGroup)

func InitRouter(e *gin.Engine) {

	// Router group
	routerGroup := e.Group(global.Cfg.Core.Prefix)

	routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Init system router
	router.InitSystemRouter(routerGroup)

	global.Logger.Sugar().Info("router register success")

}
