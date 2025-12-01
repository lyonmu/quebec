package grpc

import (
	"google.golang.org/grpc"

	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/cmd/gateway/internal/service/grpc/als"
	"github.com/lyonmu/quebec/cmd/gateway/internal/service/grpc/auth"
	"github.com/lyonmu/quebec/cmd/gateway/internal/service/grpc/proxy"
)

type Svc interface {
	Register(*grpc.Server) error
}

var registry = make([]Svc, 0)

func NewGrpcSvc(server *grpc.Server) error {

	als := als.NewAlsSvc()
	registry = append(registry, als)

	auth := auth.NewAuthSvc()
	registry = append(registry, auth)

	proxy := proxy.NewProxySvc()
	registry = append(registry, proxy)

	for _, svc := range registry {
		if err := svc.Register(server); err != nil {
			global.Logger.Sugar().Errorf("register grpc service failed: %v", err)
			return err
		}
	}
	return nil
}
