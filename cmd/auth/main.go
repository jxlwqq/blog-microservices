package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jxlwqq/blog-microservices/api/protobuf"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/interceptor"
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
	flag.Parse()
	logger := log.New()

	conf, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatal(err)
	}

	authServer, err := InitServer(logger, conf)
	if err != nil {
		logger.Fatal(err)
	}
	healthServer := health.NewServer()

	m, err := metrics.New(conf.User.Server.Name)
	if err != nil {
		logger.Fatal(err)
	}
	metricsInterceptor := interceptor.NewMetricsInterceptor(logger, m)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		metricsInterceptor.Unary(),
	)))
	protobuf.RegisterAuthServiceServer(grpcServer, authServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", conf.Auth.Server.GRPC.Port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	ch := make(chan os.Signal, 1)

	// Start gRPC server
	logger.Infof("gPRC Listening on port %s", conf.Auth.Server.GRPC.Port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Start HTTP server
	logger.Infof("HTTP Listening on port %s", conf.Auth.Server.HTTP.Port)
	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err = protobuf.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux, conf.Auth.Server.GRPC.Port, opts); err != nil {
			panic(err)
		}
		err = http.ListenAndServe(conf.Auth.Server.HTTP.Port, mux)
		if err != nil {
			panic(err)
		}
	}()

	// Start Metrics server
	logger.Infof("Metrics Listening on port %s", conf.Auth.Server.Metrics.Port)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(conf.Auth.Server.Metrics.Port, mux)
		panic(err)
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
