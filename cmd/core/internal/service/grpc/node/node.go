package node

import (
	"io"
	"time"

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
	var currentGatewayID string

	global.Logger.Sugar().Info("New Gateway connection established")

	// 循环读取流中的消息
	for {
		event, err := stream.Recv()

		// 1. 处理流断开 (EOF)
		if err == io.EOF {
			global.Logger.Sugar().Infof("Gateway %s disconnected (EOF)", currentGatewayID)
			s.handleGatewayDisconnect(currentGatewayID)
			return nil
		}

		// 2. 处理错误 (网络中断等)
		if err != nil {
			global.Logger.Sugar().Infof("Stream error from Gateway %s: %v", currentGatewayID, err)
			s.handleGatewayDisconnect(currentGatewayID)
			return err
		}

		// 3. 处理正常逻辑
		// 第一次收到消息时，锁定该流对应的 GatewayID
		if currentGatewayID == "" && event.GatewayId != "" {
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
	switch event.Event {
	case v1.EnvoyStatusEvent_CONNECT:
		global.Logger.Sugar().Infof("[CONNECT] Gateway:%s -> Node:%s", event.GatewayId, event.NodeId)
		s.registry.AddOrUpdate(event.GatewayId, &EnvoyNode{
			NodeID:    event.NodeId,
			ClusterID: event.ClusterId,
			GatewayID: event.GatewayId,
			ConnectAt: time.Now(),
		})

	case v1.EnvoyStatusEvent_DISCONNECT:
		global.Logger.Sugar().Infof("[DISCONNECT] Gateway:%s -> Node:%s", event.GatewayId, event.NodeId)
		s.registry.Remove(event.GatewayId, event.NodeId)

	case v1.EnvoyStatusEvent_HEARTBEAT:
		// 可以在这里更新 Gateway 的最后活跃时间
		global.Logger.Sugar().Infof("[HEARTBEAT] from %s", event.GatewayId)
	}
}

// handleGatewayDisconnect 当 Gateway 整个服务断开时的兜底逻辑
func (s *NodeSvc) handleGatewayDisconnect(gatewayID string) {
	if gatewayID == "" {
		return
	}
	// 策略选择：
	// 1. 悲观策略（推荐）：认为 Gateway 挂了，它管理的 Envoy 可能也都失联了，或者会重连到其他 Gateway。
	//    为了数据一致性，清理该 Gateway 的所有记录。
	s.registry.ClearGateway(gatewayID)
	global.Logger.Sugar().Infof("Cleaned up all nodes for Gateway: %s", gatewayID)

	// 2. 乐观策略：保留数据，等待 Envoy 自动重连到别的 Gateway 并触发新的 CONNECT 事件来覆盖旧数据。
	//    这种策略需要配合 Redis TTL 才能避免数据永久残留。
}
