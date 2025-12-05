package global

import (
	"github.com/lyonmu/quebec/cmd/gateway/internal/config"
	"github.com/lyonmu/quebec/pkg/tools"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	Cfg        config.Config
	Logger     *zap.Logger
	Id         tools.IDGenerator
	GrpcClient *grpc.ClientConn
)
