package main

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/comment"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
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

	commentServer, err := InitServer(logger, conf)
	if err != nil {
		logger.Fatal(err)
	}
	healthServer := health.NewServer()

	jwtManager := auth.NewJWTManager(logger, conf)

	authInterceptor := auth.NewInterceptor(logger, jwtManager, comment.AuthMethods)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))
	protobuf.RegisterCommentServiceServer(grpcServer, commentServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	logger.Info("Starting server on port " + conf.Comment.Server.Port)

	lis, err := net.Listen("tcp", conf.Comment.Server.Port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	// Start gRPC server
	ch := make(chan os.Signal, 1)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
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
