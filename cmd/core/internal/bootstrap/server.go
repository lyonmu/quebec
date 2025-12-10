package bootstrap

import (
	grpcApi "github.com/lyonmu/quebec/cmd/core/internal/api/grpc"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/cmd/core/internal/router"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

func InitServer() error {

	grpcServer, err := tools.NewGRPCServer(string(constant.ModuleNameCore), global.Metrics)
	if err != nil {
		return err
	}
	grpcApi.NewGrpcSvc(grpcServer)

	ginEngine, err := tools.NewGin(global.Metrics)
	if err != nil {
		return err
	}
	router.InitRouter(ginEngine)

	if err = tools.NewCmux(global.Cfg.Core.Port, ginEngine, grpcServer); err != nil {
		return err
	}

	return nil
}
