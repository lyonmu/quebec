package bootstrap

import (
	grpcApi "github.com/lyonmu/quebec/cmd/gateway/internal/api/grpc"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

func InitServer() error {

	conn, err := tools.NewGRPCConn(global.Cfg.Gateway.Admin, global.Metrics)
	if err != nil {
		return err
	}
	global.GrpcClient = conn

	grpcServer, err := tools.NewGRPCServer(string(constant.ModuleNameGateway), global.Metrics)
	if err != nil {
		return err
	}
	grpcApi.NewGrpcSvc(grpcServer)

	ginEngine, err := tools.NewGin(global.Metrics)
	if err != nil {
		return err
	}

	if err = tools.NewCmux(global.Cfg.Gateway.Port, ginEngine, grpcServer); err != nil {
		return err
	}

	return nil
}
