package node

import (
	"io"
	"time"

	nodepb "github.com/lyonmu/quebec/idl/node"
)

func (n *NodeSvc) NodeRegister(stream nodepb.NodeService_NodeRegisterServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		stream.Send(&nodepb.NodeRegisterReply{
			NodeId:    req.NodeId,
			Timestamp: time.Now().Unix(),
		})
	}
}
