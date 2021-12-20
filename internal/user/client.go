package user

import (
	"context"
	"github.com/jxlwqq/blog-microservices/api/protobuf"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"time"
)

func NewClient(logger *log.Logger, conf *config.Config) (protobuf.UserServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.User.Server.Host+conf.User.Server.GRPC.Port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return protobuf.NewUserServiceClient(conn), err
}
