package node

import (
	"context"
	"io"
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/ent/coregatewaycluster"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coregatewaynode"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	v1 "github.com/lyonmu/quebec/idl/node/v1"
	"google.golang.org/grpc"
)

type NodeSvc struct {
	v1.UnimplementedEnvoyRegistryServer
	registry *Registry
}

func NewNodeSvc() *NodeSvc {
	return &NodeSvc{
		registry: NewRegistry(),
	}
}

func (n *NodeSvc) Register(server *grpc.Server) error {
	global.Logger.Sugar().Info("node grpc service registered")
	v1.RegisterEnvoyRegistryServer(server, n)
	return nil
}

// SyncEnvoyStatus 处理双向流
func (s *NodeSvc) SyncEnvoyStatus(stream v1.EnvoyRegistry_SyncEnvoyStatusServer) error {
	// 用于记录当前流对应的 Gateway ID，以便断开时清理
	var currentGatewayID int64

	global.Logger.Sugar().Info("New Gateway connection established")

	// 循环读取流中的消息
	for {
		event, err := stream.Recv()

		// 1. 处理流断开 (EOF)
		if err == io.EOF {
			global.Logger.Sugar().Infof("Gateway %d disconnected (EOF)", currentGatewayID)
			s.handleGatewayDisconnect(currentGatewayID)
			return nil
		}

		// 2. 处理错误 (网络中断等)
		if err != nil {
			global.Logger.Sugar().Infof("Stream error from Gateway %d: %v", currentGatewayID, err)
			s.handleGatewayDisconnect(currentGatewayID)
			return err
		}

		// 3. 处理正常逻辑
		// 第一次收到消息时，锁定该流对应的 GatewayID
		if currentGatewayID == 0 && event.GatewayId != 0 {
			currentGatewayID = event.GatewayId
		}

		s.processEvent(event)

		// 4. (可选) 发送回执
		// 仅在需要确认或心跳时发送，避免流量风暴
		if event.Event == v1.EnvoyStatusEvent_CONNECT {
			resp := &v1.BaseResponse{
				Code:    "200",
				Message: "OK",
				Data: &v1.BaseResponse_SyncEnvoyStatus{
					SyncEnvoyStatus: &v1.SyncEnvoyStatusResponse{
						Timestamp: time.Now().Unix(),
					},
				},
			}
			if err := stream.Send(resp); err != nil {
				global.Logger.Sugar().Infof("Failed to send response: %v", err)
				return err // 发送失败通常意味着连接有问题
			}
		}
	}
}

// processEvent 处理具体的业务逻辑
func (s *NodeSvc) processEvent(event *v1.EnvoyStatusEvent) {
	requestTS := time.Now().Unix()

	switch event.Event {
	case v1.EnvoyStatusEvent_CONNECT:
		global.Logger.Sugar().Infof("[CONNECT] Gateway:%d -> Node:%s", event.GatewayId, event.NodeId)
		if err := s.registry.AddOrUpdate(event.GatewayId, &EnvoyNode{
			NodeID:        event.NodeId,
			ClusterID:     event.ClusterId,
			GatewayID:     event.GatewayId,
			ConnectAt:     time.Now(),
			LastRequestAt: requestTS,
		}); err != nil {
			global.Logger.Sugar().Errorf("upsert node failed: %v", err)
		}

	case v1.EnvoyStatusEvent_DISCONNECT:
		global.Logger.Sugar().Infof("[DISCONNECT] Gateway:%d -> Node:%s", event.GatewayId, event.NodeId)
		s.registry.Remove(event.GatewayId, event.NodeId)

	case v1.EnvoyStatusEvent_HEARTBEAT:
		// 可以在这里更新 Gateway 的最后活跃时间
		global.Logger.Sugar().Infof("[HEARTBEAT] from %d", event.GatewayId)
		if err := s.registry.UpdateLastRequestTime(event.GatewayId, event.ClusterId, event.NodeId, requestTS); err != nil {
			global.Logger.Sugar().Errorf("update last request time failed: %v", err)
		}
	}
}

// handleGatewayDisconnect 当 Gateway 整个服务断开时的兜底逻辑
func (s *NodeSvc) handleGatewayDisconnect(gatewayID int64) {
	if gatewayID == 0 {
		return
	}

	ctx := context.Background()
	now := time.Now()

	tx, err := global.EntClient.Tx(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("create tx failed: %v", err)
		return
	}

	// 统一回滚兜底
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.CoreGatewayCluster.
		Update().
		Where(
			coregatewaycluster.GatewayID(gatewayID),
			coregatewaycluster.DeletedAtIsNil(),
		).
		SetDeletedAt(now).
		Save(ctx); err != nil {
		global.Logger.Sugar().Errorf("update core_gateway_cluster failed: %v", err)
		return
	}

	if _, err = tx.CoreGatewayNode.
		Update().
		Where(
			coregatewaynode.GatewayID(gatewayID),
			coregatewaynode.DeletedAtIsNil(),
		).
		SetDeletedAt(now).
		Save(ctx); err != nil {
		global.Logger.Sugar().Errorf("update core_gateway_node failed: %v", err)
		return
	}

	if err = tx.Commit(); err != nil {
		global.Logger.Sugar().Errorf("commit tx failed: %v", err)
		return
	}

	// 事务成功后再清理内存态
	s.registry.ClearGateway(gatewayID)
	global.Logger.Sugar().Infof("Cleaned up all nodes for Gateway: %d", gatewayID)
}
