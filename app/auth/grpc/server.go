package grpc

import (
	"fmt"
	"net"

	"github.com/lyonmu/quebec/app/auth/service"
	"github.com/sirupsen/logrus"

	envoyserviceauthv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

func StartGRPCServer(port uint16) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Panicf("failed to listen to %d: %v", port, err)
	}

	users := make(map[string]string)

	users["token1"] = "user1"
	users["token2"] = "user2"
	users["token3"] = "user3"

	gs := grpc.NewServer()

	envoyserviceauthv3.RegisterAuthorizationServer(gs, service.New(users))

	logrus.Infof("starting gRPC server on: %d", port)

	if err := gs.Serve(lis); err != nil {
		logrus.Panicf("failed to start grpc server to %d: %v", port, err)
	}
}
