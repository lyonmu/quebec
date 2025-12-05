package callback

import (
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	v1pb "github.com/lyonmu/quebec/idl/node/v1"
)

func (c *XDSCallbacks) OnStreamClosed(streamID int64, node *core.Node) {
	// 原子地删除并获取值
	if value, loaded := c.sessions.LoadAndDelete(streamID); loaded {
		// 类型断言：sync.Map 存的是 interface{}，取出来需要转回 *NodeInfo
		info, ok := value.(*NodeInfo)
		if ok {
			// 发送断开事件给 Core
			c.Syncer.PushEvent(&v1pb.EnvoyStatusEvent{
				Event:     v1pb.EnvoyStatusEvent_DISCONNECT,
				NodeId:    info.NodeID,
				ClusterId: info.Cluster,
				StreamId:  info.StreamID,
			})
		}
	}
	global.Logger.Sugar().Infof("on stream closed, streamID: %d, nodeId: %s, nodeCluster: %s", streamID, node.Id, node.Cluster)
}

func (c *XDSCallbacks) OnStreamRequest(id int64, request *discoverygrpc.DiscoveryRequest) error {
	node := request.GetNode()
	global.Logger.Sugar().Infof("on stream request, streamID: %d, nodeId: %s, nodeCluster: %s", id, node.Id, node.Cluster)

	// 1. 快速路径：如果 session 已经存在，说明不是首包，直接返回，不做任何分配
	if _, ok := c.sessions.Load(id); ok {
		return nil
	}

	// 2. 组装 NodeInfo 对象
	newNodeInfo := &NodeInfo{
		NodeID:   node.Id,
		Cluster:  node.Cluster,
		StreamID: id,
	}

	// 3. 原子操作：LoadOrStore
	// 如果 loaded 为 true，说明在并发情况下，别的协程已经存进去了，我们什么都不用做
	// 如果 loaded 为 false，说明是我们存进去的，需要触发 CONNECT 事件
	_, loaded := c.sessions.LoadOrStore(id, newNodeInfo)

	if !loaded {
		// 发送连接事件给 Core
		c.Syncer.PushEvent(&v1pb.EnvoyStatusEvent{
			Event:     v1pb.EnvoyStatusEvent_CONNECT,
			NodeId:    newNodeInfo.NodeID,
			ClusterId: newNodeInfo.Cluster,
			StreamId:  newNodeInfo.StreamID,
		})
	}

	return nil
}
