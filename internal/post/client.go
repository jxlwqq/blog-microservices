package post

import (
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"google.golang.org/grpc"
)

func NewClient(postAddr string) (protobuf.PostServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(postAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := protobuf.NewPostServiceClient(conn)
	return client, conn, nil
}
