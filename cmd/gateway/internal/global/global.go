package global

import (
	"github.com/lyonmu/quebec/cmd/gateway/internal/config"
	"github.com/lyonmu/quebec/pkg/tools"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	Cfg        config.Config
	Logger     *zap.Logger
	Id         tools.IDGenerator
	GrpcClient *grpc.ClientConn
	Metrics    *prometheus.Registry
)
