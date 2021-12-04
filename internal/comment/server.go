package comment

import (
	"context"
	"fmt"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServer(repo Repository, userClient protobuf.UserServiceClient, postClient protobuf.PostServiceClient) protobuf.CommentServiceServer {
	return &Server{
		repo:       repo,
		userClient: userClient,
		postClient: postClient,
	}
}

type Server struct {
	protobuf.UnimplementedCommentServiceServer
	repo       Repository
	userClient protobuf.UserServiceClient
	postClient protobuf.PostServiceClient
}

func (s Server) CreateComment(ctx context.Context, req *protobuf.CreateCommentRequest) (*protobuf.CreateCommentResponse, error) {
	userID := ctx.Value("ID").(uint64)
	fmt.Println("userID: ", userID)
	userResp, err := s.userClient.GetUser(ctx, &protobuf.GetUserRequest{Id: userID})
	user := userResp.GetUser()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found: %v", userID, err)
	}
	postID := req.GetComment().GetPostId()
	postResp, err := s.postClient.GetPost(ctx, &protobuf.GetPostRequest{Id: postID})
	post := postResp.GetPost()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found: %v", postID, err)
	}
	comment := Comment{
		Content: req.GetComment().GetContent(),
		PostID:  postID,
		UserID:  userID,
	}
	err = s.repo.Create(&comment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create comment: %v", err)
	}
	resp := &protobuf.CreateCommentResponse{
		Comment: &protobuf.ResponseComment{
			Id:        comment.ID,
			Content:   comment.Content,
			PostId:    comment.PostID,
			UserId:    comment.UserID,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
			User:      user,
			Post:      post,
		},
	}

	return resp, nil
}

func (s Server) UpdateComment(ctx context.Context, req *protobuf.UpdateCommentRequest) (*protobuf.UpdateCommentResponse, error) {
	panic("implement me")
}

func (s Server) DeleteComment(ctx context.Context, req *protobuf.DeleteCommentRequest) (*protobuf.DeleteCommentResponse, error) {
	panic("implement me")
}

func (s Server) GetCommentsByPostID(ctx context.Context, req *protobuf.GetCommentListByPostIDRequest) (*protobuf.GetCommentListByPostIDResponse, error) {
	panic("implement me")
}
