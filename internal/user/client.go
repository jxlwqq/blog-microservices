package user

import (
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"google.golang.org/grpc"
)

func NewClient(userAddr string) (protobuf.UserServiceClient, error) {
	conn, err := grpc.Dial(userAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return protobuf.NewUserServiceClient(conn), err
}
