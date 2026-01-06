package als

import (
	"fmt"
	"io"

	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/pkg/mq/kafka"

	accesslogv3 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type AlsSvc struct{}

func NewAlsSvc() *AlsSvc {
	return &AlsSvc{}
}

func (a *AlsSvc) Register(server *grpc.Server) error {
	accesslogv3.RegisterAccessLogServiceServer(server, a)
	global.Logger.Sugar().Info("als service registered")
	return nil
}

func (a *AlsSvc) StreamAccessLogs(stream accesslogv3.AccessLogService_StreamAccessLogsServer) error {
	ctx := stream.Context()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			global.Logger.Sugar().Info("als event stream closed by envoy")
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("als event recv error: %v", err)
		}

		// 发送 proto 二进制数据到 Kafka
		if global.AlsKafkaProducer != nil {
			data, err := proto.Marshal(in)
			if err != nil {
				global.Logger.Sugar().Errorf("als event marshal error: %v", err)
				continue
			}
			producer := global.AlsKafkaProducer.(*kafka.Producer[[]byte, []byte])
			if err := producer.Produce(ctx, nil, &data); err != nil {
				global.Logger.Sugar().Errorf("send als to kafka error: %v", err)
			}
		}

		// 保留日志输出用于调试
		global.Logger.Sugar().Debug("als event received")
	}
}
