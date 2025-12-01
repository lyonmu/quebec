package proxy

import (
	"context"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
)

// defaultCallbacks 实现默认的 callback 来记录日志信息
type defaultCallbacks struct{}

func (c *defaultCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	global.Logger.Sugar().Infof("on stream opened, streamID: %d, typeURL: %s", streamID, typeURL)
	return nil
}

func (c *defaultCallbacks) OnStreamClosed(streamID int64, node *core.Node) {
	global.Logger.Sugar().Infof("on stream closed, streamID: %d, node: %+v", streamID, node)
}

func (c *defaultCallbacks) OnStreamRequest(streamID int64, request *discoverygrpc.DiscoveryRequest) error {
	global.Logger.Sugar().Infof("on stream request, streamID: %d, request: %+v", streamID, request)
	return nil
}

func (c *defaultCallbacks) OnStreamResponse(ctx context.Context, streamID int64, request *discoverygrpc.DiscoveryRequest, response *discoverygrpc.DiscoveryResponse) {
	global.Logger.Sugar().Infof("on stream response, streamID: %d, request: %+v, response: %+v", streamID, request, response)
}

func (c *defaultCallbacks) OnFetchRequest(ctx context.Context, request *discoverygrpc.DiscoveryRequest) error {
	global.Logger.Sugar().Infof("on fetch request, request: %+v", request)
	return nil
}

func (c *defaultCallbacks) OnFetchResponse(request *discoverygrpc.DiscoveryRequest, response *discoverygrpc.DiscoveryResponse) {
	global.Logger.Sugar().Infof("on fetch response, request: %+v, response: %+v", request, response)
}

// Delta xDS 相关方法
func (c *defaultCallbacks) OnDeltaStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	global.Logger.Sugar().Infof("on delta stream opened, streamID: %d, typeURL: %s", streamID, typeURL)
	return nil
}

func (c *defaultCallbacks) OnDeltaStreamClosed(streamID int64, node *core.Node) {
	global.Logger.Sugar().Infof("on delta stream closed, streamID: %d, node: %+v", streamID, node)
}

func (c *defaultCallbacks) OnStreamDeltaRequest(streamID int64, request *discoverygrpc.DeltaDiscoveryRequest) error {
	global.Logger.Sugar().Infof("on stream delta request, streamID: %d, request: %+v", streamID, request)
	return nil
}

func (c *defaultCallbacks) OnStreamDeltaResponse(streamID int64, request *discoverygrpc.DeltaDiscoveryRequest, response *discoverygrpc.DeltaDiscoveryResponse) {
	global.Logger.Sugar().Infof("on stream delta response_nonce, streamID: %d, request: %+v, response: %+v", streamID, request, response)
}
