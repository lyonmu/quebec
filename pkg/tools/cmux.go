package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/pkg/metrics"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewGin() (*gin.Engine, error) {

	r := gin.New()
	r.Use(cors.Default())
	r.Use(gin.Recovery())
	if err := metrics.RegisterMetrics(r); err != nil {
		return nil, err
	}
	if gin.Mode() != gin.ReleaseMode {
		r.Use(gin.Logger())
		pprof.Register(r)
	}
	return r, nil
}

func NewGRPCServer(name string, opts ...grpc.ServerOption) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(opts...)
	healthServer := health.NewServer()
	healthServer.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	return grpcServer, nil
}

func NewCmux(p uint16, r *gin.Engine, g *grpc.Server) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
	if err != nil {
		return err
	}
	m := cmux.New(l)
	go g.Serve(m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc")))
	go (&http.Server{Handler: r}).Serve(m.Match(cmux.HTTP1Fast()))

	return m.Serve()
}
