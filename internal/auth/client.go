package auth

import (
	"context"
	"time"

	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/auth/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(logger log.Logger, conf *config.Config) (v1.AuthServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		conf.Auth.Server.Host+conf.Auth.Server.GRPC.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclogging.UnaryClientInterceptor(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	client := v1.NewAuthServiceClient(conn)
	return client, nil
}
