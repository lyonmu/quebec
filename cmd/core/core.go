package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/lyonmu/quebec/cmd/core/internal/bootstrap"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/common"
	"github.com/prometheus/common/version"
)

// swag init -g core.go -o ./internal/docs --parseDependency --parseInternal

// @title                       quebec core Swagger API接口文档
// @description                 quebec core 模块后端
// @version                     v0.0.1
// @license.name                MIT
// @license.url                 https://github.com/lyonmu/quebec/blob/master/LICENSE
// @contact.name                Lyon Mu
// @contact.url                 https://github.com/lyonmu
// @contact.email               lyonmu@foxmail.com
// @host                        http://localhost:59024
// @BasePath                    /core/api
// @schemes                     http
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 token认证

func main() {
	kong.Parse(&global.Cfg,
		kong.Name(string(common.ModuleNameCore)),
		kong.Description(string(common.ModuleNameCore)),
		kong.UsageOnError(),
		kong.HelpOptions{
			Compact: true,
			Summary: true,
		},
	)
	if global.Cfg.Version {
		fmt.Println(version.Print(string(common.ModuleNameCore)))
		os.Exit(0)
	}
	bootstrap.Start()
}
