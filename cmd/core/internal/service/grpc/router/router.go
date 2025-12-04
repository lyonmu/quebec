package router

import (
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"google.golang.org/grpc"
)

type RouterSvc struct{}

func NewRouterSvc() *RouterSvc {
	return &RouterSvc{}
}

func (r *RouterSvc) Register(server *grpc.Server) error {

	global.Logger.Sugar().Info("router grpc service registered")
	return nil
}
