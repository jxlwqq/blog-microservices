package user

import (
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"google.golang.org/grpc"
)

func NewClient(conf *config.Config) (protobuf.UserServiceClient, error) {
	conn, err := grpc.Dial(conf.User.Server.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return protobuf.NewUserServiceClient(conn), err
}
