package grpc

import (
	"google.golang.org/grpc"

	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/cmd/core/internal/service/grpc/node"
	"github.com/lyonmu/quebec/cmd/core/internal/service/grpc/router"
)

type Svc interface {
	Register(*grpc.Server) error
}

var registry = make([]Svc, 0)

func NewGrpcSvc(server *grpc.Server) error {

	node := node.NewNodeSvc()
	registry = append(registry, node)

	router := router.NewRouterSvc()
	registry = append(registry, router)

	for _, svc := range registry {
		if err := svc.Register(server); err != nil {
			global.Logger.Sugar().Errorf("register grpc service failed: %v", err)
			return err
		}
	}
	return nil
}
