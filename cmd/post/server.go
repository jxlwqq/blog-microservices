package main

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/pkg/dbcontext"
	"github.com/stonecutter/blog-microservices/internal/post"
	"github.com/stonecutter/blog-microservices/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "config file")

func main() {
	flag.Parse()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	db, err := dbcontext.New(conf.Post.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	postRepo := post.NewRepository(db)
	userClient, userConn, err := user.NewClient(conf.User.Server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	postServer := post.NewServer(postRepo, userClient)

	jwtManager := auth.NewJWTManager(conf.JWT.Secret, conf.JWT.Expires)

	methods := make(map[string]bool)
	prefix := "/api.protobuf.PostService/"
	methods[prefix+"CreatePost"] = true // 需要jwt验证
	methods[prefix+"UpdatePost"] = true
	methods[prefix+"DeletePost"] = true
	methods[prefix+"GetPost"] = false // 不需要jwt验证
	methods[prefix+"ListPost"] = false

	interceptor := auth.NewInterceptor(jwtManager, methods)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
	)
	protobuf.RegisterPostServiceServer(grpcServer, postServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	log.Println("Starting server on port " + conf.Post.Server.Port)

	lis, err := net.Listen("tcp", conf.Post.Server.Port)

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
	select {
	case <-ctx.Done():
		close(ch)
	}
	fmt.Println("Graceful Shutdown end")
}
