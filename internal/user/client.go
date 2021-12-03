package user

import (
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"google.golang.org/grpc"
)

func NewClient(userAddr string) (protobuf.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(userAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return protobuf.NewUserServiceClient(conn), conn, err
}
