package callback

import (
	"context"
	"sync"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/cmd/gateway/internal/service/grpc/node"
	"google.golang.org/grpc"
)

type NodeInfo struct {
	NodeID   string
	Cluster  string
	StreamID int64
}

type XDSCallbacks struct {
	Syncer   *node.CoreSyncer
	sessions sync.Map
}

// NewGatewayCallbacks 初始化 CoreSyncer 并构建 xDS 回调
// conn: 到 Core 服务的 gRPC 连接
// gatewayID: 当前 Gateway 的唯一标识
func NewGatewayCallbacks(conn *grpc.ClientConn, gatewayID string) server.Callbacks {
	// 1. 初始化 CoreSyncer
	// 注意：这里我们假设 NewCoreSyncer 已经在同一个包或引用的包中定义
	syncer := node.NewCoreSyncer(conn, gatewayID)

	// 2. 启动后台同步协程
	// 这一步封装在这里，外部调用者就不需要关心内部有个 syncer 需要启动了
	syncer.Start()

	global.Logger.Sugar().Infof("CoreSyncer started for GatewayID: %s", gatewayID)

	// 3. 返回组装好的 Callbacks
	// sync.Map 的零值即直接可用，无需显式初始化
	return &XDSCallbacks{
		Syncer: syncer,
	}
}

func (c *XDSCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	global.Logger.Sugar().Infof("on stream opened, streamID: %d, typeURL: %s", streamID, typeURL)
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
