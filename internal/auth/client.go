package auth

import (
	"context"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/auth/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"time"
)

func NewClient(logger *log.Logger, conf *config.Config) (v1.AuthServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.Auth.Server.Host+conf.Auth.Server.GRPC.Port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := v1.NewAuthServiceClient(conn)
	return client, nil
}
