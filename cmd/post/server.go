package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/pkg/interceptor"
	"github.com/stonecutter/blog-microservices/internal/pkg/jwt"
	"github.com/stonecutter/blog-microservices/internal/pkg/log"
	"github.com/stonecutter/blog-microservices/internal/pkg/metrics"
	"github.com/stonecutter/blog-microservices/internal/post"
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

	jwtManager := jwt.NewJWTManager(logger, conf)

	authInterceptor := interceptor.NewAuthInterceptor(logger, jwtManager, post.AuthMethods)
	m, err := metrics.New(conf.Post.Server.Name)
	if err != nil {
		logger.Fatal(err)
	}
	metricsInterceptor := interceptor.NewMetricsInterceptor(logger, m)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		authInterceptor.Unary(),
		metricsInterceptor.Unary(),
	)))

	protobuf.RegisterPostServiceServer(grpcServer, postServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", conf.Post.Server.Port)

	// Start gRPC server
	ch := make(chan os.Signal, 1)
	logger.Infof("gPRC Listening on port %s", conf.Post.Server.Port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	logger.Infof("HTTP Listening on port %s", conf.Post.HTTP.Port)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err = http.ListenAndServe(conf.Post.HTTP.Port, nil); err != nil {
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
