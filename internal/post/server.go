package post

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/stonecutter/blog-microservices/api/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServer(repo Repository, userClient protobuf.UserServiceClient) protobuf.PostServiceServer {
	return &Server{repo: repo, userClient: userClient}
}

type Server struct {
	protobuf.UnimplementedPostServiceServer
	repo       Repository
	userClient protobuf.UserServiceClient
}

func (s Server) GetPost(ctx context.Context, req *protobuf.GetPostRequest) (*protobuf.GetPostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post, err := s.repo.Get(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found: %v", err)
	}
	user, err := s.userClient.GetUser(ctx, &protobuf.GetUserRequest{Id: post.UserID})
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	protobufPost := entityToProtobuf(post, user.User)
	resp := &protobuf.GetPostResponse{
		Post: protobufPost,
	}

	return resp, nil
}

func (s Server) CreatePost(ctx context.Context, req *protobuf.CreatePostRequest) (*protobuf.CreatePostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, ok := ctx.Value("ID").(uint64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	post := &Post{
		Title:   req.GetPost().GetTitle(),
		Content: req.GetPost().GetContent(),
		UserID:  userID,
	}
	err = s.repo.Create(post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}

	user, err := s.userClient.GetUser(ctx, &protobuf.GetUserRequest{Id: post.UserID})
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	resp := &protobuf.CreatePostResponse{
		Post: entityToProtobuf(post, user.User),
	}

	return resp, nil
}

func (s Server) UpdatePost(ctx context.Context, req *protobuf.UpdatePostRequest) (*protobuf.UpdatePostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post := &Post{
		ID:      req.GetId(),
		Title:   req.GetPost().GetTitle(),
		Content: req.GetPost().GetContent(),
	}
	userID, ok := ctx.Value("ID").(uint64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	if post.UserID != userID {
		return nil, status.Errorf(codes.Unauthenticated, "user %d is not the owner of post %d", userID, post.ID)
	}
	err = s.repo.Update(post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post: %v", err)
	}

	resp := &protobuf.UpdatePostResponse{
		Success: true,
	}

	return resp, nil
}

func (s Server) DeletePost(ctx context.Context, req *protobuf.DeletePostRequest) (*protobuf.DeletePostResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post, err := s.repo.Get(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found", req.GetId())
	}
	userID, ok := ctx.Value("ID").(uint64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	if post.UserID != userID {
		return nil, status.Errorf(codes.Unauthenticated, "user %d is not the owner of post %d", userID, post.ID)
	}

	err = s.repo.Delete(post.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}

	resp := &protobuf.DeletePostResponse{
		Success: true,
	}

	return resp, nil
}

func (s Server) ListPosts(ctx context.Context, req *protobuf.ListPostsRequest) (*protobuf.ListPostsResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	list, err := s.repo.List(int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list posts: %v", err)
	}

	var userIDs []uint64

	for _, post := range list {
		userIDs = append(userIDs, post.UserID)
	}
	userReq := &protobuf.GetUserListByIDsRequest{
		Ids: userIDs,
	}

	userResp, err := s.userClient.GetUserListByIDs(ctx, userReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	users := userResp.GetUsers()
	var posts []*protobuf.ResponsePost
	for _, post := range list {
		user := &protobuf.User{}
		for _, item := range users {
			if post.UserID == item.Id {
				_ = copier.Copy(user, item)
			}
		}
		posts = append(posts, entityToProtobuf(post, user))
	}

	count, err := s.repo.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count posts: %v", err)
	}

	resp := &protobuf.ListPostsResponse{
		Posts: posts,
		Count: count,
	}
	return resp, nil
}

func entityToProtobuf(post *Post, user *protobuf.User) *protobuf.ResponsePost {
	return &protobuf.ResponsePost{
		Id:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserId:    post.UserID,
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
		User:      user,
	}
}
