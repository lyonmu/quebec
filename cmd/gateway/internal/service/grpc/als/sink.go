package als

import (
	"fmt"
	"io"

	"github.com/lyonmu/quebec/cmd/gateway/internal/global"

	accesslogv3 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type AlsSvc struct {
	marshaler protojson.MarshalOptions
}

func NewAlsSvc() *AlsSvc {
	return &AlsSvc{
		marshaler: protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   true,
		},
	}
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

		b, err := a.marshaler.Marshal(in)
		if err != nil {
			return fmt.Errorf("als event marshal error: %v", err)
		}

		global.Logger.Sugar().Infof("als event: %s", string(b))
	}
}
