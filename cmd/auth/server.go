package main

import (
	flag "github.com/spf13/pflag"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	conn, err := grpc.Dial(conf.User.Server.Host+conf.User.Server.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	userClient := protobuf.NewUserServiceClient(conn)

	jwtManager := auth.NewJWTManager(conf.JWT.Secret, conf.JWT.Expires)
	authServer := auth.NewServer(userClient, jwtManager)
	protobuf.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", conf.Auth.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on port %s", conf.Auth.Server.Port)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
