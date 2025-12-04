package xds

import (
	"context"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	xdsv3cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xdsv3log "github.com/envoyproxy/go-control-plane/pkg/log"
	xdsv3server "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/cmd/gateway/internal/service/grpc/xds/callback"
	"google.golang.org/grpc"
)

var (
	xdsCache = xdsv3cache.NewSnapshotCache(false, xdsv3cache.IDHash{}, &xdsv3log.DefaultLogger{})
)

type ProxySvc struct {
	xdsserver xdsv3server.Server
}

// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/bootstrap/v3/bootstrap.proto#config-bootstrap-v3-bootstrap-dynamicresources
func (s *ProxySvc) Register(gs *grpc.Server) error {
	// 注册所有需要的 XDS 服务
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(gs, s.xdsserver) // ADS 必须
	clusterservice.RegisterClusterDiscoveryServiceServer(gs, s.xdsserver)   // CDS
	endpointservice.RegisterEndpointDiscoveryServiceServer(gs, s.xdsserver) // EDS
	routeservice.RegisterRouteDiscoveryServiceServer(gs, s.xdsserver)       // RDS
	listenerservice.RegisterListenerDiscoveryServiceServer(gs, s.xdsserver) // LDS
	global.Logger.Info("xDS services registered successfully")
	return nil
}

func NewProxySvc() *ProxySvc {

	// create default callback instance to record envoy xDS logs
	callbacks := &callback.XDSCallbacks{}

	return &ProxySvc{xdsv3server.NewServer(context.Background(), xdsCache, callbacks)}
}
