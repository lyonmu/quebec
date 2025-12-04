package node

import (
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"google.golang.org/grpc"
)

type NodeSvc struct {
}

func NewNodeSvc() *NodeSvc {
	return &NodeSvc{}
}

func (n *NodeSvc) Register(server *grpc.Server) error {
	global.Logger.Sugar().Info("node grpc service registered")
	return nil
}
