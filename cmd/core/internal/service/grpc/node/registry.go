package node

import (
	"context"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coregatewaycluster"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coregatewaynode"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
)

// EnvoyNode 表示单个 Envoy 实例的信息
type EnvoyNode struct {
	NodeID    string
	ClusterID string
	Address   string
	ConnectAt time.Time
	// LastRequestAt keeps the latest request/heartbeat timestamp reported by Envoy.
	LastRequestAt int64
	GatewayID     int64 // 记录是通过哪个 Gateway 连接的
}

// Registry 是一个线程安全的存储库
type Registry struct {
	// 使用互斥锁保证并发下的顺序化写入
	mu sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{}
}

// AddOrUpdate 注册/更新节点
func (r *Registry) AddOrUpdate(gatewayID int64, node *EnvoyNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lastRequest := node.LastRequestAt
	if lastRequest == 0 {
		lastRequest = time.Now().Unix()
	}

	if err := r.upsertCluster(ctx, node.ClusterID, gatewayID, lastRequest); err != nil {
		return err
	}

	if err := r.upsertNode(ctx, gatewayID, node, lastRequest); err != nil {
		return err
	}

	return nil
}

// Remove 移除节点
func (r *Registry) Remove(gatewayID int64, nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := global.EntClient.CoreGatewayNode.Update().
		Where(
			coregatewaynode.NodeID(nodeID),
			coregatewaynode.GatewayID(gatewayID),
			coregatewaynode.DeletedAtIsNil(),
		).
		SetDeletedAt(time.Now()).
		Save(ctx); err != nil {
		global.Logger.Sugar().Errorf("soft delete core_gateway_node failed: %s", err)
	}
}

// ClearGateway 清理指定 Gateway 的所有数据 (用于 Gateway 断连时)
func (r *Registry) ClearGateway(gatewayID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := global.EntClient.CoreGatewayNode.Update().
		Where(
			coregatewaynode.GatewayID(gatewayID),
			coregatewaynode.DeletedAtIsNil(),
		).
		SetDeletedAt(time.Now()).
		Save(ctx); err != nil {
		global.Logger.Sugar().Errorf("soft delete core_gateway_node by gateway failed: %s", err)
	}
}

// GetAll 获取所有在线节点 (用于调试或监控)
func (r *Registry) GetAll() []*EnvoyNode {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	nodes, err := global.EntClient.CoreGatewayNode.Query().
		Where(coregatewaynode.DeletedAtIsNil()).
		All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("query core_gateway_node failed: %s", err)
		return nil
	}

	list := make([]*EnvoyNode, 0, len(nodes))
	for _, n := range nodes {
		list = append(list, &EnvoyNode{
			NodeID:        n.NodeID,
			ClusterID:     n.ClusterID,
			GatewayID:     n.GatewayID,
			LastRequestAt: n.NodeLastRequestTime,
			ConnectAt:     n.CreatedAt,
		})
	}
	return list
}

// UpdateLastRequestTime 更新节点及所在集群的最新请求时间
func (r *Registry) UpdateLastRequestTime(gatewayID int64, clusterID, nodeID string, ts int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	requestTS := ts
	if requestTS == 0 {
		requestTS = time.Now().Unix()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := r.upsertCluster(ctx, clusterID, gatewayID, requestTS); err != nil {
		return err
	}

	if _, err := global.EntClient.CoreGatewayNode.Update().
		Where(
			coregatewaynode.NodeID(nodeID),
			coregatewaynode.ClusterID(clusterID),
			coregatewaynode.GatewayID(gatewayID),
			coregatewaynode.DeletedAtIsNil(),
		).
		SetNodeLastRequestTime(requestTS).
		Save(ctx); err != nil {
		global.Logger.Sugar().Errorf("update core_gateway_node last_request_time failed: %s", err)
		return err
	}

	return nil
}

func (r *Registry) upsertCluster(ctx context.Context, clusterID string, gatewayID int64, lastRequest int64) error {
	if clusterID == "" {
		return nil
	}

	// upsert cluster record and bump last_request_time
	err := global.EntClient.CoreGatewayCluster.Create().
		SetClusterID(clusterID).
		SetGatewayID(gatewayID).
		SetClusterCreateTime(time.Now().Unix()).
		OnConflict(
			sql.ConflictColumns(coregatewaycluster.FieldClusterID),
		).
		Update(func(u *ent.CoreGatewayClusterUpsert) {
			u.SetGatewayID(gatewayID)
			u.SetClusterLastRequestTime(lastRequest)
			u.ClearDeletedAt()
			u.SetUpdatedAt(time.Now())
		}).
		Exec(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("upsert core_gateway_cluster failed: %s", err)
		return err
	}
	return nil
}

func (r *Registry) upsertNode(ctx context.Context, gatewayID int64, node *EnvoyNode, lastRequest int64) error {
	err := global.EntClient.CoreGatewayNode.Create().
		SetNodeID(node.NodeID).
		SetClusterID(node.ClusterID).
		SetGatewayID(gatewayID).
		SetNodeRegisterTime(time.Now().Unix()).
		SetNodeLastRequestTime(lastRequest).
		OnConflict(
			sql.ConflictColumns(coregatewaynode.FieldNodeID),
		).
		Update(func(u *ent.CoreGatewayNodeUpsert) {
			u.SetClusterID(node.ClusterID)
			u.SetGatewayID(gatewayID)
			u.SetNodeLastRequestTime(lastRequest)
			u.ClearDeletedAt()
			u.SetUpdatedAt(time.Now())
		}).
		Exec(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("upsert core_gateway_node failed: %s", err)
		return err
	}
	return nil
}
