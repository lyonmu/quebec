package bootstrap

import (
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/cmd/core/internal/router"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

func InitServer() error {

	grpcServer, err := tools.NewGRPCServer(string(constant.ModuleNameCore))
	if err != nil {
		return err
	}

	ginEngine, err := tools.NewGin()
	if err != nil {
		return err
	}
	router.InitRouter(ginEngine)

	if err = tools.NewCmux(global.Cfg.Core.Port, ginEngine, grpcServer); err != nil {
		return err
	}

	return nil
}
