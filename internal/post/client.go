package post

import (
	"context"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"google.golang.org/grpc"
	"time"
)

func NewClient(conf *config.Config) (protobuf.PostServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.Post.Server.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := protobuf.NewPostServiceClient(conn)
	return client, nil
}
