package router

import (
	"github.com/gin-gonic/gin"
	v1Api "github.com/lyonmu/quebec/cmd/core/internal/api/http/v1"
	_ "github.com/lyonmu/quebec/cmd/core/internal/docs"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/cmd/core/internal/middleware/http"
	v1Route "github.com/lyonmu/quebec/cmd/core/internal/router/v1"
	"github.com/lyonmu/quebec/cmd/core/internal/web"
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
	group := e.Group(global.Cfg.Core.Prefix)

	// swagger
	swaggerRouter := group.Group("swagger")
	{

		if gin.Mode() != gin.ReleaseMode {
			swaggerRouter.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		} else {
			swaggerRouter.Use(http.JwtAuth()).GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	// Init system router
	v1route.InitSystemRouter(group, v1api)

	if err := web.Register(e, "/"); err != nil {
		global.Logger.Sugar().Warnf("register embedded web failed: %v", err)
	}

	global.Logger.Sugar().Info("router http register success")

}
