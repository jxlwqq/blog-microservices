package post

import (
	"context"
	"github.com/jxlwqq/blog-microservices/api/protobuf/post/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"time"
)

func NewClient(logger *log.Logger, conf *config.Config) (v1.PostServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.Post.Server.Host+conf.Post.Server.GRPC.Port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := v1.NewPostServiceClient(conn)
	return client, nil
}
