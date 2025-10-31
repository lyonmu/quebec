package common

import (
	cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
)

type LoadBalancerPolicy int8

// 负载均衡策略 1:随机,2:加权轮询,3:加权最小请求数,4:环形一致性哈希,5:Maglev一致性哈希
// 参考: https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/load_balancing/load_balancers#arch-overview-load-balancing-types
const (
	LoadBalancerPolicyRandom       LoadBalancerPolicy = 1
	LoadBalancerPolicyRoundRobin   LoadBalancerPolicy = 2
	LoadBalancerPolicyLeastRequest LoadBalancerPolicy = 3
	LoadBalancerPolicyRingHash     LoadBalancerPolicy = 4
	LoadBalancerPolicyMaglev       LoadBalancerPolicy = 5
)

var LoadBalancerPolicyMap = map[LoadBalancerPolicy]cluster_v3.Cluster_LbPolicy{
	1: cluster_v3.Cluster_RANDOM,        // 随机
	2: cluster_v3.Cluster_ROUND_ROBIN,   // 加权轮询
	3: cluster_v3.Cluster_LEAST_REQUEST, // 加权最小请求数
	4: cluster_v3.Cluster_RING_HASH,     // 环形一致性哈希
	5: cluster_v3.Cluster_MAGLEV,        // Maglev一致性哈希
}
