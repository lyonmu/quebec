package xds

import (
	"fmt"
	"os"
	"strings"
	"time"

	accesslogv3 "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	accessloggrpcv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/grpc/v3"
	extauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"

	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/google/uuid"
	"github.com/lyonmu/quebec/cmd/gateway/internal/common"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Create a cache

type SvcInfo struct {
	Name      string
	Prefix    string
	Instances []InstanceInfo
}

type InstanceInfo struct {
	Id   string
	Host string
	Port uint32
}

func (s *SvcInfo) MakeEndpoint() *endpoint.ClusterLoadAssignment {
	hs := make([]*endpoint.LocalityLbEndpoints, 0, len(s.Instances))
	for _, ip := range s.Instances {
		// 直接使用 IP 地址
		hs = append(hs, &endpoint.LocalityLbEndpoints{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  ip.Host,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: ip.Port,
									},
								},
							},
						},
					},
				},
			}},
		})
	}
	return &endpoint.ClusterLoadAssignment{
		ClusterName: s.Name,
		Endpoints:   hs,
	}
}

func (s *SvcInfo) MakeCluster() *cluster.Cluster {
	return &cluster.Cluster{
		Name:                 s.Name,
		ConnectTimeout:       durationpb.New(3 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
		LbPolicy:             common.LBPolicyMap[common.LoadBalancerPolicyMaglev],
		EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
			EdsConfig: &core.ConfigSource{
				ResourceApiVersion: core.ApiVersion_V3,
				ConfigSourceSpecifier: &core.ConfigSource_Ads{
					Ads: &core.AggregatedConfigSource{},
				},
			},
			ServiceName: s.Name,
		},
		DnsLookupFamily: cluster.Cluster_AUTO,
		// 添加 DNS 解析相关配置
		RespectDnsTtl:  true,
		DnsRefreshRate: durationpb.New(30 * time.Second),
	}
}

func (s *SvcInfo) MakeRoute() *route.Route {

	routeConfig := &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{Prefix: fmt.Sprintf("%s%s", global.Cfg.Gateway.Prefix, s.Prefix)},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: s.Name,
				},
				PrefixRewrite: s.Prefix,
			},
		},
	}

	// 如果开启权限认证，则添加 ext_authz 配置
	// if global.Cfg.AuthConfig.Enabled {
	// 如果需要权限校验，添加 ext_authz 配置
	// if !s.NoAuth {
	// 	routeConfig.TypedPerFilterConfig = map[string]*anypb.Any{
	// 		"envoy.filters.http.ext_authz": createExtAuthzGrpcConfig(),
	// 	}
	// }
	// }

	return routeConfig
}

// createExtAuthzGrpcConfig 创建 gRPC ext_authz 配置
func createExtAuthzGrpcConfig() *anypb.Any {
	extAuthzConfig := &extauthz.ExtAuthz{
		Services: &extauthz.ExtAuthz_GrpcService{
			GrpcService: &core.GrpcService{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
						ClusterName: common.ClusterName,
					},
				},
				Timeout: durationpb.New(60 * time.Second),
			},
		},
		FailureModeAllow: false, // 认证失败时拒绝请求
		WithRequestBody: &extauthz.BufferSettings{
			MaxRequestBytes:     4096,
			AllowPartialMessage: true,
		},
		ClearRouteCache: false,
		// 设置 gRPC 服务的方法
		TransportApiVersion: core.ApiVersion_V3,
	}

	config, _ := anypb.New(extAuthzConfig)
	return config
}

// RouteConfiguration
func MakeRouteConfig(svcs []*SvcInfo) *route.RouteConfiguration {

	rs := make([]*route.Route, 0, len(svcs))
	for _, s := range svcs {
		rs = append(rs, s.MakeRoute())
	}

	return &route.RouteConfiguration{
		Name: common.RouteName,
		VirtualHosts: []*route.VirtualHost{
			{
				Name:    common.VirtualHostName,
				Domains: []string{"*"}, // 只允许一个 "*"
				Routes:  rs,
			},
		},
	}
}

// Listener
func MakeListener(routeCfg *route.RouteConfiguration) *listener.Listener {
	routerConfig, _ := anypb.New(&router.Router{})

	// 创建 HTTP 过滤器列表
	httpFilters := []*hcm.HttpFilter{}

	// 如果启用认证，先添加 ext_authz 过滤器
	// if global.Cfg.Gateway.EnableAuth {
	// 	httpFilters = append(httpFilters, &hcm.HttpFilter{
	// 		Name: "envoy.filters.http.ext_authz",
	// 		ConfigType: &hcm.HttpFilter_TypedConfig{
	// 			TypedConfig: createExtAuthzGrpcConfig(),
	// 		},
	// 	})
	// }

	// router 过滤器必须是最后一个
	httpFilters = append(httpFilters, &hcm.HttpFilter{
		Name: common.HttpFilterName, // router 过滤器
		ConfigType: &hcm.HttpFilter_TypedConfig{
			TypedConfig: routerConfig,
		},
	})

	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: common.HttpStatPrefixName,
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				RouteConfigName: common.RouteName, // 对应 snapshot 的 key
				ConfigSource: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{},
				},
			},
		},
		HttpFilters: httpFilters,
		// HTTP 请求优化配置
		RequestTimeout:      durationpb.New(60 * time.Second),    // 减少请求超时时间
		StreamIdleTimeout:   durationpb.New(60 * time.Second),    // 减少流空闲超时
		MaxRequestHeadersKb: &wrapperspb.UInt32Value{Value: 256}, // 限制请求头大小
	}

	// 根据 debug 日志级别动态添加 gRPC Access Log
	if strings.ToLower(global.Cfg.Log.Level) == "debug" {
		accessLogConfig, err := anypb.New(&accessloggrpcv3.HttpGrpcAccessLogConfig{
			CommonConfig: &accessloggrpcv3.CommonGrpcAccessLogConfig{
				LogName:             common.AccessLogName,
				TransportApiVersion: core.ApiVersion_V3,
				GrpcService: &core.GrpcService{
					TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
							ClusterName: common.ClusterName,
						},
					},
					Timeout: durationpb.New(60 * time.Second),
				},
			},
		})
		if err != nil {
			global.Logger.Sugar().Errorf("failed to marshal HttpGrpcAccessLogConfig: %v", err)
		} else {
			manager.AccessLog = []*accesslogv3.AccessLog{
				{
					Name:       common.AccessLogName,
					ConfigType: &accesslogv3.AccessLog_TypedConfig{TypedConfig: accessLogConfig},
				},
			}
		}
	}

	pbst, err := anypb.New(manager)
	if err != nil {
		global.Logger.Sugar().Errorf("failed to marshal HttpConnectionManager: %v", err)
		os.Exit(1)
	}

	return &listener.Listener{
		Name: common.ListenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(global.Cfg.Gateway.Port),
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{
			{
				Filters: []*listener.Filter{
					{
						Name:       common.ListenerFilterName,
						ConfigType: &listener.Filter_TypedConfig{TypedConfig: pbst},
					},
				},
			},
		},
	}
}

// Snapshot
func GenerateSnapshot(svcs []*SvcInfo) *cache.Snapshot {
	// 1. 生成集群配置 CDS
	clusters := make([]types.Resource, 0, len(svcs))
	for _, s := range svcs {
		cluster := s.MakeCluster()
		clusters = append(clusters, cluster)
	}

	// 2. 生成端点配置 EDS
	endpoints := make([]types.Resource, 0, len(svcs))
	for _, s := range svcs {
		endpoint := s.MakeEndpoint()
		endpoints = append(endpoints, endpoint)
	}

	// 3. 生成路由配置 RDS
	routeCfg := MakeRouteConfig(svcs)

	// 4. 生成监听器配置 LDS
	listenerCfg := MakeListener(routeCfg)
	// 5. 创建 snapshot
	resources := map[resource.Type][]types.Resource{
		resource.ClusterType:  clusters,
		resource.EndpointType: endpoints,
		resource.RouteType:    {routeCfg},
		resource.ListenerType: {listenerCfg},
	}

	snap, err := cache.NewSnapshot(uuid.Must(uuid.NewV7()).String(), resources)
	if err != nil {
		global.Logger.Sugar().Errorf("failed to generate snapshot: %v", err)
		os.Exit(1)
	}

	// 验证 snapshot
	if err := snap.Consistent(); err != nil {
		global.Logger.Sugar().Errorf("snapshot validation failed: %v", err)
		os.Exit(1)
	}

	global.Logger.Sugar().Info("snapshot generated and validated successfully")
	return snap
}
