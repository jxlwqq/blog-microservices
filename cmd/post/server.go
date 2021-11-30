package main

import (
	flag "github.com/spf13/pflag"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/config"
	"github.com/stonecutter/blog-microservices/internal/post"
	"github.com/stonecutter/blog-microservices/internal/user"
	"github.com/stonecutter/blog-microservices/pkg/dbcontext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "config file")

func main() {
	flag.Parse()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	db, err := dbcontext.NewDB(conf.Post.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	postRepo := post.NewRepository(db)
	userClient, err := user.NewClient(conf.User.Server.Host + conf.User.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	postServer := post.NewServer(postRepo, userClient)

	jwtManager := auth.NewJWTManager(conf.JWT.Secret, conf.JWT.Expires)

	methods := make(map[string]bool)

	methods["/api.protobuf.PostService/CreatePost"] = true // 需要jwt验证
	methods["/api.protobuf.PostService/UpdatePost"] = true
	methods["/api.protobuf.PostService/DeletePost"] = true
	methods["/api.protobuf.PostService/GetPost"] = false // 不需要jwt验证
	methods["/api.protobuf.PostService/ListPost"] = false

	interceptor := auth.NewInterceptor(jwtManager, methods)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
	)
	protobuf.RegisterPostServiceServer(grpcServer, postServer)
	reflection.Register(grpcServer)

	log.Println("Starting server on port " + conf.Post.Server.Port)

	lis, err := net.Listen("tcp", conf.Post.Server.Port)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
