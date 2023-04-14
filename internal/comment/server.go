package comment

import (
	"context"

	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/comment/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func NewServer(logger log.Logger, repo Repository) v1.CommentServiceServer {
	return &Server{
		logger: logger,
		repo:   repo,
	}
}

type Server struct {
	v1.UnimplementedCommentServiceServer
	logger log.Logger
	repo   Repository
}

func (s Server) CreateComment(ctx context.Context, req *v1.CreateCommentRequest) (*v1.CreateCommentResponse, error) {

	find, err := s.repo.GetByUUID(ctx, req.GetComment().GetUuid())
	if err == nil {
		return &v1.CreateCommentResponse{
			Comment: entityToProtobuf(find),
		}, nil
	}

	comment := &Comment{
		UUID:    req.GetComment().GetUuid(),
		Content: req.GetComment().GetContent(),
		PostID:  req.GetComment().GetPostId(),
		UserID:  req.GetComment().GetUserId(),
	}
	err = s.repo.Create(ctx, comment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create comment: %v", err)
	}

	return &v1.CreateCommentResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}

func (s Server) CreateCommentCompensate(ctx context.Context, req *v1.CreateCommentRequest) (*v1.CreateCommentResponse, error) {
	err := s.repo.DeleteByUUID(ctx, req.GetComment().GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create comment: %v", err)
	}

	return &v1.CreateCommentResponse{}, nil
}

func (s Server) GetComment(ctx context.Context, req *v1.GetCommentRequest) (*v1.GetCommentResponse, error) {
	id := req.GetId()
	comment, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "could not get comment: %v", err)
	}

	return &v1.GetCommentResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}

func (s Server) GetCommentByUUID(ctx context.Context, req *v1.GetCommentByUUIDRequest) (*v1.GetCommentByUUIDResponse, error) {
	comment, err := s.repo.GetByUUID(ctx, req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get comment: %v", err)
	}

	return &v1.GetCommentByUUIDResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}

func (s Server) UpdateComment(ctx context.Context, req *v1.UpdateCommentRequest) (*v1.UpdateCommentResponse, error) {

	commentID := req.GetComment().GetId()
	comment, err := s.repo.Get(ctx, commentID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "comment %d not found: %v", commentID, err)
	}

	if req.GetComment().GetContent() != "" {
		comment.Content = req.GetComment().GetContent()
	}

	err = s.repo.Update(ctx, comment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not update comment: %v", err)
	}

	return &v1.UpdateCommentResponse{Success: true}, nil

}

func (s Server) DeleteComment(ctx context.Context, req *v1.DeleteCommentRequest) (*v1.DeleteCommentResponse, error) {
	commentID := req.GetId()
	_, err := s.repo.Get(ctx, commentID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "comment %d not found: %v", commentID, err)
	}
	err = s.repo.Delete(ctx, commentID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete comment: %v", err)
	}
	return &v1.DeleteCommentResponse{
		Success: true,
	}, nil
}

func (s Server) DeleteCommentCompensate(ctx context.Context, req *v1.DeleteCommentRequest) (*v1.DeleteCommentResponse, error) {
	commentID := req.GetId()
	comment, err := s.repo.GetWithUnscoped(ctx, commentID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "comment %d not found: %v", commentID, err)
	}
	comment.DeletedAt = gorm.DeletedAt{}
	err = s.repo.UpdateWithUnscoped(ctx, comment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete comment: %v", err)
	}
	return &v1.DeleteCommentResponse{
		Success: true,
	}, nil
}

func (s Server) DeleteCommentsByPostID(ctx context.Context, req *v1.DeleteCommentsByPostIDRequest) (*v1.DeleteCommentsByPostIDResponse, error) {
	postID := req.GetPostId()
	err := s.repo.DeleteByPostID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete comments: %v", err)
	}
	return &v1.DeleteCommentsByPostIDResponse{
		Success: true,
	}, nil
}

func (s Server) DeleteCommentsByPostIDCompensate(ctx context.Context, req *v1.DeleteCommentsByPostIDRequest) (*v1.DeleteCommentsByPostIDResponse, error) {
	postID := req.GetPostId()

	err := s.repo.UpdateByPostIDWithUnscoped(ctx, postID, Comment{DeletedAt: gorm.DeletedAt{}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete comments: %v", err)
	}
	return &v1.DeleteCommentsByPostIDResponse{
		Success: true,
	}, nil
}

func (s Server) ListCommentsByPostID(ctx context.Context, req *v1.ListCommentsByPostIDRequest) (*v1.ListCommentsByPostIDResponse, error) {
	postID := req.GetPostId()
	offset := req.GetOffset()
	limit := req.GetOffset()
	if limit == 0 {
		limit = 10 // default limit
	}
	list, err := s.repo.ListByPostID(ctx, postID, int(offset), int(limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get comments: %v", err)
	}

	var comments []*v1.Comment

	for _, comment := range list {
		comments = append(comments, entityToProtobuf(comment))
	}

	total, err := s.repo.CountByPostID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get comments: %v", err)
	}

	return &v1.ListCommentsByPostIDResponse{
		Comments: comments,
		Total:    total,
	}, nil
}

func entityToProtobuf(comment *Comment) *v1.Comment {
	return &v1.Comment{
		Id:        comment.ID,
		Uuid:      comment.UUID,
		Content:   comment.Content,
		PostId:    comment.PostID,
		UserId:    comment.UserID,
		CreatedAt: timestamppb.New(comment.CreatedAt),
		UpdatedAt: timestamppb.New(comment.UpdatedAt),
	}
}
