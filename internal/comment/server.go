package comment

import (
	"context"
	"github.com/jxlwqq/blog-microservices/api/protobuf"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var prefix = "/api.protobuf.CommentService/"

var AuthMethods = map[string]bool{
	prefix + "CreateComment":          true,
	prefix + "UpdateComment":          true,
	prefix + "DeleteComment":          true,
	prefix + "GetCommentListByPostID": false,
}

func NewServer(logger *log.Logger, repo Repository, userClient protobuf.UserServiceClient, postClient protobuf.PostServiceClient) protobuf.CommentServiceServer {
	return &Server{
		logger:     logger,
		repo:       repo,
		userClient: userClient,
		postClient: postClient,
	}
}

type Server struct {
	protobuf.UnimplementedCommentServiceServer
	logger     *log.Logger
	repo       Repository
	userClient protobuf.UserServiceClient
	postClient protobuf.PostServiceClient
}

func (s Server) CreateComment(ctx context.Context, req *protobuf.CreateCommentRequest) (*protobuf.CreateCommentResponse, error) {
	userID, ok := ctx.Value("ID").(uint64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	s.logger.Info("CreateComment", "userID", userID)
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
		Comment: &protobuf.Comment{
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
	userID, ok := ctx.Value("ID").(uint64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	s.logger.Info("UpdateComment", "userID", userID)
	userResp, err := s.userClient.GetUser(ctx, &protobuf.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found: %v", userID, err)
	}
	user := userResp.GetUser()
	commentID := req.GetComment().GetId()
	c, err := s.repo.Get(commentID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "comment %d not found: %v", commentID, err)
	}
	if c.UserID != user.GetId() {
		return nil, status.Errorf(codes.PermissionDenied, "user %d does not own comment %d", userID, commentID)
	}

	comment := &Comment{
		ID:      commentID,
		Content: req.GetComment().GetContent(),
	}

	err = s.repo.Update(comment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not update comment: %v", err)
	}

	return &protobuf.UpdateCommentResponse{Success: true}, nil

}

func (s Server) DeleteComment(ctx context.Context, req *protobuf.DeleteCommentRequest) (*protobuf.DeleteCommentResponse, error) {
	userID, ok := ctx.Value("ID").(uint64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	userResp, err := s.userClient.GetUser(ctx, &protobuf.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found: %v", userID, err)
	}
	user := userResp.GetUser()
	commentID := req.GetId()
	comment, err := s.repo.Get(commentID)
	postResp, err := s.postClient.GetPost(ctx, &protobuf.GetPostRequest{Id: comment.PostID})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found: %v", comment.PostID, err)
	}
	post := postResp.GetPost()

	if comment.UserID != user.GetId() && post.UserId != user.GetId() {
		return nil, status.Errorf(codes.PermissionDenied, "user %d does not own comment %d", userID, commentID)
	}
	err = s.repo.Delete(commentID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete comment: %v", err)
	}

	return &protobuf.DeleteCommentResponse{
		Success: true,
	}, nil
}

func (s Server) GetCommentsByPostID(ctx context.Context, req *protobuf.GetCommentListByPostIDRequest) (*protobuf.GetCommentListByPostIDResponse, error) {
	postID := req.GetPostId()
	comments, err := s.repo.ListByPostID(postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get comments: %v", err)
	}

	var userIDs []uint64

	for _, comment := range comments {
		userIDs = append(userIDs, comment.UserID)
	}

	usersResp, err := s.userClient.GetUserListByIDs(ctx, &protobuf.GetUserListByIDsRequest{Ids: userIDs})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get users: %v", err)
	}

	users := usersResp.GetUsers()

	var commentsWithUsers []*protobuf.Comment

	for _, comment := range comments {
		for _, user := range users {
			if user.GetId() == comment.UserID {
				commentsWithUsers = append(commentsWithUsers, &protobuf.Comment{
					Id:        comment.ID,
					Content:   comment.Content,
					PostId:    comment.PostID,
					UserId:    comment.UserID,
					CreatedAt: timestamppb.New(comment.CreatedAt),
					UpdatedAt: timestamppb.New(comment.UpdatedAt),
					User:      user,
				})
			}
		}
	}

	return &protobuf.GetCommentListByPostIDResponse{
		Comments: commentsWithUsers,
	}, nil

}
