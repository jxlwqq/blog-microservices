package main

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/blog/v1"
	"github.com/jxlwqq/blog-microservices/internal/blog"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/interceptor"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/pkg/metrics"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "path to config file")

func main() {
	logger := log.New()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatal("load config failed", err)
	}

	blogServer, err := InitServer(logger, conf)
	if err != nil {
		logger.Fatal("init server failed", err)
	}

	healthServer := health.NewServer()

	jwtManager := jwt.NewManager(logger, conf)
	authInterceptor := interceptor.NewAuthInterceptor(logger, jwtManager, blog.AuthMethods)
	m, err := metrics.New(conf.Blog.Server.Name)
	if err != nil {
		logger.Fatal("init metrics failed", err)
	}
	metricsInterceptor := interceptor.NewMetricsInterceptor(logger, m)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		authInterceptor.Unary(),
		metricsInterceptor.Unary(),
		grpc_recovery.UnaryServerInterceptor(),
	)))

	v1.RegisterBlogServiceServer(grpcServer, blogServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	listen, err := net.Listen("tcp", conf.Blog.Server.GRPC.Port)
	if err != nil {
		logger.Fatal("listen failed", err)
	}

	err = grpcServer.Serve(listen)
	if err != nil {
		logger.Fatal("grpc serve failed", err)
	}
}
