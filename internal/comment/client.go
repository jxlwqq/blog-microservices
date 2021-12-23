package comment

import (
	"context"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/comment/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"time"
)

func NewClient(logger *log.Logger, conf *config.Config) (v1.CommentServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.Comment.Server.Host+conf.Comment.Server.GRPC.Port, grpc.WithInsecure())
	if err != nil {
		logger.Error("grpc dial error", err)
		return nil, err
	}
	return v1.NewCommentServiceClient(conn), nil
}
