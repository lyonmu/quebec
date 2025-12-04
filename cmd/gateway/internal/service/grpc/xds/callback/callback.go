package callback

import (
	"context"
	"sync"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/cmd/gateway/internal/service/grpc/node"
)

type NodeInfo struct {
	NodeID  string
	Cluster string
	Address string
}

type XDSCallbacks struct {
	Syncer *node.CoreSyncer
	// 存储 StreamID -> NodeID/Cluster 的映射，用于断开时查找
	mu       sync.Mutex
	sessions map[int64]*NodeInfo
}

func (c *XDSCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	global.Logger.Sugar().Infof("on stream opened, streamID: %d, typeURL: %s", streamID, typeURL)
	return nil
}

func (c *XDSCallbacks) OnStreamClosed(streamID int64, node *core.Node) {
	global.Logger.Sugar().Infof("on stream closed, streamID: %d, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		streamID, node.Id, node.Cluster, node.Metadata)
}

func (c *XDSCallbacks) OnStreamRequest(streamID int64, request *discoverygrpc.DiscoveryRequest) error {
	global.Logger.Sugar().Infof("on stream request, streamID: %d, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		streamID, request.Node.Id, request.Node.Cluster, request.Node.Metadata)
	return nil
}

func (c *XDSCallbacks) OnStreamResponse(ctx context.Context, streamID int64, request *discoverygrpc.DiscoveryRequest, response *discoverygrpc.DiscoveryResponse) {
	global.Logger.Sugar().Infof("on stream response, streamID: %d, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		streamID, request.Node.Id, request.Node.Cluster, request.Node.Metadata)
}

func (c *XDSCallbacks) OnFetchRequest(ctx context.Context, request *discoverygrpc.DiscoveryRequest) error {
	global.Logger.Sugar().Infof("on fetch request, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		request.Node.Id, request.Node.Cluster, request.Node.Metadata)
	return nil
}

func (c *XDSCallbacks) OnFetchResponse(request *discoverygrpc.DiscoveryRequest, response *discoverygrpc.DiscoveryResponse) {
	global.Logger.Sugar().Infof("on fetch response, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		request.Node.Id, request.Node.Cluster, request.Node.Metadata)
}

// Delta xDS 相关方法
func (c *XDSCallbacks) OnDeltaStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	global.Logger.Sugar().Infof("on delta stream opened, streamID: %d, typeURL: %s", streamID, typeURL)
	return nil
}

func (c *XDSCallbacks) OnDeltaStreamClosed(streamID int64, node *core.Node) {
	global.Logger.Sugar().Infof("on delta stream closed, streamID: %d, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		streamID, node.Id, node.Cluster, node.Metadata)
}

func (c *XDSCallbacks) OnStreamDeltaRequest(streamID int64, request *discoverygrpc.DeltaDiscoveryRequest) error {
	global.Logger.Sugar().Infof("on stream delta request, streamID: %d, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		streamID, request.Node.Id, request.Node.Cluster, request.Node.Metadata)
	return nil
}

func (c *XDSCallbacks) OnStreamDeltaResponse(streamID int64, request *discoverygrpc.DeltaDiscoveryRequest, response *discoverygrpc.DeltaDiscoveryResponse) {
	global.Logger.Sugar().Infof("on stream delta response_nonce, streamID: %d, nodeId: %s, nodeCluster: %s, node.Metadata: %+v",
		streamID, request.Node.Id, request.Node.Cluster, request.Node.Metadata)
}
