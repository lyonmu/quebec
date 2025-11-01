package common

import (
	cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
)

type LBPolicy int8

// 负载均衡策略 1:随机,2:加权轮询,3:加权最小请求数,4:环形一致性哈希,5:Maglev一致性哈希
// 参考: https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/load_balancing/load_balancers#arch-overview-load-balancing-types
const (
	LoadBalancerPolicyRandom       LBPolicy = 1
	LoadBalancerPolicyRoundRobin   LBPolicy = 2
	LoadBalancerPolicyLeastRequest LBPolicy = 3
	LoadBalancerPolicyRingHash     LBPolicy = 4
	LoadBalancerPolicyMaglev       LBPolicy = 5
)

var LBPolicyMap = map[LBPolicy]cluster_v3.Cluster_LbPolicy{
	1: cluster_v3.Cluster_RANDOM,        // 随机
	2: cluster_v3.Cluster_ROUND_ROBIN,   // 加权轮询
	3: cluster_v3.Cluster_LEAST_REQUEST, // 加权最小请求数
	4: cluster_v3.Cluster_RING_HASH,     // 环形一致性哈希
	5: cluster_v3.Cluster_MAGLEV,        // Maglev一致性哈希
}

const (
	DataPlane          = "quebec_gateway_data_plane"
	ControlPlane       = "quebec_gateway_control_plane"
	ClusterName        = "quebec_gateway_cluster"
	RouteName          = "quebec_gateway_route"
	ListenerName       = "quebec_gateway_listener"
	ListenerFilterName = "quebec_gateway_listener_filter"
	HttpStatPrefixName = "quebec_gateway_http"
	HttpFilterName     = "quebec_gateway_http_filter"
	VirtualHostName    = "quebec_gateway_virtual_host"
)
