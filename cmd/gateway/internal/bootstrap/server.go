package bootstrap

import (
	grpcApi "github.com/lyonmu/quebec/cmd/gateway/internal/api/grpc"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

func InitServer() error {

	grpcServer, err := tools.NewGRPCServer(string(constant.ModuleNameGateway))
	if err != nil {
		return err
	}
	grpcApi.NewGrpcSvc(grpcServer)

	ginEngine, err := tools.NewGin()
	if err != nil {
		return err
	}

	if err = tools.NewCmux(global.Cfg.Gateway.Port, ginEngine, grpcServer); err != nil {
		return err
	}

	return nil
}
