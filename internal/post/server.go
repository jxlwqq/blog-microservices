package post

import (
	"context"
	"github.com/jxlwqq/blog-microservices/api/protobuf/post/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServer(logger *log.Logger, repo Repository) v1.PostServiceServer {
	return &Server{logger: logger, repo: repo}
}

type Server struct {
	v1.UnimplementedPostServiceServer
	logger *log.Logger
	repo   Repository
}

func (s Server) IncrementCommentsCount(ctx context.Context, req *v1.IncrementCommentsCountRequest) (*v1.IncrementCommentsCountResponse, error) {
	postID := req.GetId()
	p, err := s.repo.Get(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found", postID)
	}
	p.CommentsCount++

	err = s.repo.Update(ctx, p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post %d", postID)
	}

	return &v1.IncrementCommentsCountResponse{Success: true}, nil
}

func (s Server) IncrementCommentsCountCompensate(ctx context.Context, req *v1.IncrementCommentsCountRequest) (*v1.IncrementCommentsCountResponse, error) {
	postID := req.GetId()
	p, err := s.repo.Get(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found", postID)
	}
	p.CommentsCount--
	err = s.repo.Update(ctx, p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post %d", postID)
	}
	return &v1.IncrementCommentsCountResponse{Success: true}, nil
}

func (s Server) DecrementCommentsCount(ctx context.Context, request *v1.DecrementCommentsCountRequest) (*v1.DecrementCommentsCountResponse, error) {
	postID := request.GetId()
	p, err := s.repo.Get(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found", postID)
	}
	p.CommentsCount--
	err = s.repo.Update(ctx, p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post %d", postID)
	}
	return &v1.DecrementCommentsCountResponse{Success: true}, nil
}

func (s Server) GetPost(ctx context.Context, req *v1.GetPostRequest) (*v1.GetPostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found: %v", err)
	}
	protobufPost := entityToProtobuf(post)
	resp := &v1.GetPostResponse{
		Post: protobufPost,
	}

	return resp, nil
}

func (s Server) CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	post := &Post{
		UUID:    req.GetPost().GetUuid(),
		Title:   req.GetPost().GetTitle(),
		Content: req.GetPost().GetContent(),
		UserID:  req.GetPost().GetUserId(),
	}
	err = s.repo.Create(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}

	resp := &v1.CreatePostResponse{
		Post: entityToProtobuf(post),
	}

	return resp, nil
}

func (s Server) UpdatePost(ctx context.Context, req *v1.UpdatePostRequest) (*v1.UpdatePostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	postID := req.GetPost().GetId()
	post, err := s.repo.Get(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found", postID)
	}

	if req.GetPost().GetTitle() != "" {
		post.Title = req.GetPost().GetTitle()
	}

	if req.GetPost().GetContent() != "" {
		post.Content = req.GetPost().GetContent()
	}

	s.logger.Info("update post", post)

	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post: %v", err)
	}

	resp := &v1.UpdatePostResponse{
		Success: true,
	}

	return resp, nil
}

func (s Server) DeletePost(ctx context.Context, req *v1.DeletePostRequest) (*v1.DeletePostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found", req.GetId())
	}

	err = s.repo.Delete(ctx, post.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}

	resp := &v1.DeletePostResponse{
		Success: true,
	}

	return resp, nil
}

func (s Server) ListPosts(ctx context.Context, req *v1.ListPostsRequest) (*v1.ListPostsResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	list, err := s.repo.List(ctx, int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list posts: %v", err)
	}

	var posts []*v1.Post
	for _, post := range list {
		posts = append(posts, entityToProtobuf(post))
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count posts: %v", err)
	}

	resp := &v1.ListPostsResponse{
		Posts: posts,
		Count: count,
	}
	return resp, nil
}

func entityToProtobuf(post *Post) *v1.Post {
	return &v1.Post{
		Id:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		CommentsCount: post.CommentsCount,
		UserId:        post.UserID,
		CreatedAt:     timestamppb.New(post.CreatedAt),
		UpdatedAt:     timestamppb.New(post.UpdatedAt),
	}
}
