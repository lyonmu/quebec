package node

import (
	"sync"
	"time"
)

// EnvoyNode 表示单个 Envoy 实例的信息
type EnvoyNode struct {
	NodeID    string
	ClusterID string
	Address   string
	ConnectAt time.Time
	GatewayID string // 记录是通过哪个 Gateway 连接的
}

// Registry 是一个线程安全的存储库
type Registry struct {
	// 使用 sync.Map 或者 RWMutex 保护 map
	// 结构建议： map[GatewayID]map[NodeID]*EnvoyNode
	// 这样当 Gateway 断开时，可以快速清理该 Gateway 下的所有节点
	mu    sync.RWMutex
	store map[string]map[string]*EnvoyNode
}

func NewRegistry() *Registry {
	return &Registry{
		store: make(map[string]map[string]*EnvoyNode),
	}
}

// AddOrUpdate 注册/更新节点
func (r *Registry) AddOrUpdate(gatewayID string, node *EnvoyNode) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[gatewayID]; !ok {
		r.store[gatewayID] = make(map[string]*EnvoyNode)
	}
	r.store[gatewayID][node.NodeID] = node
}

// Remove 移除节点
func (r *Registry) Remove(gatewayID string, nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if nodes, ok := r.store[gatewayID]; ok {
		delete(nodes, nodeID)
		// 如果该 Gateway 下没节点了，可选清理 map
		if len(nodes) == 0 {
			delete(r.store, gatewayID)
		}
	}
}

// ClearGateway 清理指定 Gateway 的所有数据 (用于 Gateway 断连时)
func (r *Registry) ClearGateway(gatewayID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, gatewayID)
}

// GetAll 获取所有在线节点 (用于调试或监控)
func (r *Registry) GetAll() []*EnvoyNode {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*EnvoyNode, 0)
	for _, nodes := range r.store {
		for _, node := range nodes {
			list = append(list, node)
		}
	}
	return list
}
