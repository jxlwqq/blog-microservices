package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/blog/v1"
	"github.com/jxlwqq/blog-microservices/internal/blog"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/interceptor"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	ch := make(chan os.Signal, 1)

	// Start gRPC server
	logger.Infof("gRPC Listening on port %s", conf.Blog.Server.GRPC.Port)
	go func() {
		err = grpcServer.Serve(listen)
		if err != nil {
			logger.Fatal("grpc serve failed", err)
		}
	}()

	// Start HTTP server
	logger.Infof("HTTP Listening on port %s", conf.Blog.Server.HTTP.Port)
	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err = v1.RegisterBlogServiceHandlerFromEndpoint(context.Background(), mux, conf.Blog.Server.GRPC.Port, opts); err != nil {
			logger.Fatal(err)
		}
		err = http.ListenAndServe(conf.Blog.Server.HTTP.Port, mux)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Start Metrics server
	logger.Infof("Metrics Listening on port %s", conf.Blog.Server.Metrics.Port)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(conf.Blog.Server.Metrics.Port, mux)
		panic(err)
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	<-ctx.Done()
	close(ch)
	logger.Info("Graceful Shutdown end")
}
