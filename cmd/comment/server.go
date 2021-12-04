package main

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/comment"
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

	db, err := dbcontext.NewDB(conf.Comment.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	repo := comment.NewRepository(db)

	userClient, userConn, err := user.NewClient(conf.User.Server.Host + conf.User.Server.Port)
	defer userConn.Close()
	if err != nil {
		log.Fatal(err)
	}

	postClient, postConn, err := post.NewClient(conf.Post.Server.Host + conf.Post.Server.Port)
	defer postConn.Close()
	if err != nil {
		log.Fatal(err)
	}

	commentServer := comment.NewServer(repo, userClient, postClient)

	jwtManager := auth.NewJWTManager(conf.JWT.Secret, conf.JWT.Expires)

	prefix := "/api.protobuf.CommentService/"
	methods := map[string]bool{
		prefix + "CreateComment":          true,
		prefix + "UpdateComment":          true,
		prefix + "DeleteComment":          true,
		prefix + "GetCommentListByPostID": false,
	}

	authInterceptor := auth.NewInterceptor(jwtManager, methods)

	healthServer := health.NewServer()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))

	protobuf.RegisterCommentServiceServer(grpcServer, commentServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	log.Println("Starting server on port " + conf.Comment.Server.Port)

	lis, err := net.Listen("tcp", conf.Comment.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
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
	select {
	case <-ctx.Done():
		close(ch)
	}
	fmt.Println("Graceful Shutdown end")

}
