package post

import (
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"google.golang.org/grpc"
)

func NewClient(conf *config.Config) (protobuf.PostServiceClient, error) {
	conn, err := grpc.Dial(conf.Post.Server.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := protobuf.NewPostServiceClient(conn)
	return client, nil
}
