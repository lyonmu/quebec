package tools

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	grpc_prome "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/lyonmu/quebec/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewGin(reg *prometheus.Registry) (*gin.Engine, error) {

	r := gin.New()
	r.Use(cors.Default())
	r.Use(gin.Recovery())
	if err := metrics.RegisterMetrics(r, reg); err != nil {
		return nil, err
	}
	if gin.Mode() != gin.ReleaseMode {
		r.Use(gin.Logger())
		pprof.Register(r)
	}
	return r, nil
}

func NewGRPCServer(name string, reg *prometheus.Registry, opts ...grpc.ServerOption) (*grpc.Server, error) {

	if len(name) == 0 {
		return nil, fmt.Errorf("server name cannot be empty")
	}

	if reg == nil {
		return nil, fmt.Errorf("prometheus registry cannot be nil")
	}

	counterOpts := []grpc_prome.CounterOption{
		grpc_prome.WithNamespace(name),
	}
	histOpts := []grpc_prome.HistogramOption{
		grpc_prome.WithHistogramNamespace(name),
	}
	s := grpc_prome.NewServerMetrics(grpc_prome.WithServerCounterOptions(counterOpts...), grpc_prome.WithServerHandlingTimeHistogram(histOpts...))

	defaultOpts := []grpc.ServerOption{
		grpc.StreamInterceptor(s.StreamServerInterceptor()),
		grpc.UnaryInterceptor(s.UnaryServerInterceptor()),
	}

	// Merge default options with custom options (custom options take precedence)
	allOpts := make([]grpc.ServerOption, 0, len(defaultOpts)+len(opts))
	allOpts = append(allOpts, defaultOpts...)
	allOpts = append(allOpts, opts...)
	reg.MustRegister(s)

	grpcServer := grpc.NewServer(allOpts...)
	healthServer := health.NewServer()
	healthServer.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	return grpcServer, nil
}

// NewGRPCConn creates a new gRPC client connection with sensible defaults.
// The caller is responsible for closing the connection when no longer needed.
func NewGRPCConn(endpoint string, reg *prometheus.Registry, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	if len(endpoint) == 0 {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	if reg == nil {
		return nil, fmt.Errorf("prometheus registry cannot be nil")
	}
	c := grpc_prome.NewClientMetrics()

	// Base options that are always applied
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(c.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(c.StreamClientInterceptor()),
	}

	// Merge default options with custom options (custom options take precedence)
	allOpts := make([]grpc.DialOption, 0, len(defaultOpts)+len(opts))
	allOpts = append(allOpts, defaultOpts...)
	allOpts = append(allOpts, opts...)
	reg.MustRegister(c)

	// Create the connection
	conn, err := grpc.NewClient(endpoint, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", endpoint, err)
	}

	return conn, nil
}

func NewCmux(p uint16, r *gin.Engine, g *grpc.Server) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
	if err != nil {
		return err
	}
	mux := cmux.New(l)
	// 匹配 gRPC (HTTP/2 + gRPC content-type) - 最精确的匹配
	grpcL := mux.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"),
	)

	// 匹配标准 HTTP/1.1 - 使用完整解析的 HTTP1 匹配器
	httpL := mux.Match(cmux.HTTP1())

	// 兜底匹配器（可选，用于调试或拒绝未知协议）
	anyL := mux.Match(cmux.Any()) // 匹配所有剩余流量

	// 启动 gRPC 服务器
	go func() {
		if err := g.Serve(grpcL); err != nil && !errors.Is(err, cmux.ErrListenerClosed) {
			log.Panicf("[ERROR] gRPC server failed: %v", err)
		}
	}()

	// 启动 HTTP 服务器
	httpServer := &http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := httpServer.Serve(httpL); err != nil && !errors.Is(err, cmux.ErrListenerClosed) {
			log.Panicf("[ERROR] HTTP server failed: %v", err)
		}
	}()

	// 处理未知协议的连接（安全加固）
	go func() {
		for {
			conn, err := anyL.Accept()
			if err != nil {
				if errors.Is(err, cmux.ErrListenerClosed) {
					return
				}
				continue
			}
			// 立即关闭未知协议连接
			go func(c net.Conn) {
				defer c.Close()
				// 可选：发送简单的拒绝消息
				c.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nUnsupported protocol\r\n"))
			}(conn)
		}
	}()

	return mux.Serve()
}
