package node

import (
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	nodepb "github.com/lyonmu/quebec/idl/node"
	"google.golang.org/grpc"
)

type NodeSvc struct {
	nodepb.UnimplementedNodeServiceServer
}

func NewNodeSvc() *NodeSvc {
	return &NodeSvc{}
}

func (n *NodeSvc) Register(server *grpc.Server) error {
	nodepb.RegisterNodeServiceServer(server, n)
	global.Logger.Sugar().Info("node service registered")
	return nil
}
