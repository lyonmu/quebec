package node

import (
	"context"
	"log"
	"time"

	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	v1 "github.com/lyonmu/quebec/idl/node/v1"

	"google.golang.org/grpc"
)

// CoreSyncer 负责与 Core 服务保持长连接并推送事件
type CoreSyncer struct {
	client    v1.EnvoyRegistryClient
	sendCh    chan *v1.EnvoyStatusEvent
	gatewayID string
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewCoreSyncer(conn *grpc.ClientConn, gatewayID string) *CoreSyncer {
	ctx, cancel := context.WithCancel(context.Background())
	return &CoreSyncer{
		client:    v1.NewEnvoyRegistryClient(conn),
		sendCh:    make(chan *v1.EnvoyStatusEvent, 1000), // 缓冲区设大一点
		gatewayID: gatewayID,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start 启动后台发送循环
func (s *CoreSyncer) Start() {
	go s.runLoop()
}

func (s *CoreSyncer) runLoop() {
	var stream v1.EnvoyRegistry_SyncEnvoyStatusClient
	var err error

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			// 1. 建立连接
			if stream == nil {
				global.Logger.Sugar().Info("Connecting to core ...")
				stream, err = s.client.SyncEnvoyStatus(s.ctx)
				if err != nil {
					global.Logger.Sugar().Errorf("Failed to connect core : %v, retrying in 5s...", err)
					time.Sleep(5 * time.Second)
					continue
				}
			}

			// 2. 监听事件并发送
			select {
			case event := <-s.sendCh:
				event.GatewayId = s.gatewayID // 自动填充 Gateway ID
				if err := stream.Send(event); err != nil {
					global.Logger.Sugar().Errorf("Stream send error: %v, reconnecting...", err)
					stream = nil // 标记重连
				}
			case <-s.ctx.Done():
				return
			}
		}
	}
}

// PushEvent 供 xDS 回调函数调用，非阻塞
func (s *CoreSyncer) PushEvent(event *v1.EnvoyStatusEvent) {
	select {
	case s.sendCh <- event:
	default:
		log.Println("CoreSyncer channel full, dropping event for node:", event.NodeId)
	}
}
