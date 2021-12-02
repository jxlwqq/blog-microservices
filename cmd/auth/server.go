package main

import (
	flag "github.com/spf13/pflag"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/config"
	"github.com/stonecutter/blog-microservices/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "path to config file")

func main() {
	flag.Parse()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()

	userClient, err := user.NewClient(conf.User.Server.Host + conf.User.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	jwtManager := auth.NewJWTManager(conf.JWT.Secret, conf.JWT.Expires)
	authServer := auth.NewServer(userClient, jwtManager)
	protobuf.RegisterAuthServiceServer(grpcServer, authServer)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", conf.Auth.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on port %s", conf.Auth.Server.Port)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
