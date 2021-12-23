package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jxlwqq/blog-microservices/api/protobuf/post/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/interceptor"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "config file")

func main() {
	flag.Parse()
	logger := log.New()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatal(err)
	}

	postServer, err := InitServer(logger, conf)
	if err != nil {
		logger.Fatal(err)
	}
	healthServer := health.NewServer()

	m, err := metrics.New(conf.Post.Server.Name)
	if err != nil {
		logger.Fatal(err)
	}
	metricsInterceptor := interceptor.NewMetricsInterceptor(logger, m)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		metricsInterceptor.Unary(),
		grpc_recovery.UnaryServerInterceptor(),
	)))

	v1.RegisterPostServiceServer(grpcServer, postServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", conf.Post.Server.GRPC.Port)

	// Start gRPC server
	ch := make(chan os.Signal, 1)
	logger.Infof("gPRC Listening on port %s", conf.Post.Server.GRPC.Port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// todo Start HTTP server

	// Start Metrics server
	logger.Infof("Metrics Listening on port %s", conf.Post.Server.Metrics.Port)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err = http.ListenAndServe(conf.Post.Server.Metrics.Port, nil); err != nil {
			panic(err)
		}
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	<-ctx.Done()
	close(ch)
	fmt.Println("Graceful Shutdown end")
}
