package post

import (
	"context"

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

func (s Server) GetPost(ctx context.Context, request *protobuf.GetPostRequest) (*protobuf.GetPostResponse, error) {
	post, err := s.repo.Get(request.GetId())
	if err != nil {
		return nil, err
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
	ID := ctx.Value("ID").(uint64)
	post := &Post{
		Title:   req.GetPost().GetTitle(),
		Content: req.GetPost().GetContent(),
		UserID:  ID,
	}
	err := s.repo.Create(post)
	if err != nil {
		return nil, err
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

func (s Server) UpdatePost(ctx context.Context, request *protobuf.UpdatePostRequest) (*protobuf.UpdatePostResponse, error) {
	post := &Post{
		ID:      request.GetId(),
		Title:   request.GetPost().GetTitle(),
		Content: request.GetPost().GetContent(),
	}
	ID := ctx.Value("ID").(uint64)
	if post.UserID != ID {
		return nil, status.Errorf(codes.Unauthenticated, "user %d is not the owner of post %d", ID, post.ID)
	}
	err := s.repo.Update(post)
	if err != nil {
		return nil, err
	}

	resp := &protobuf.UpdatePostResponse{
		Success: true,
	}

	return resp, nil
}

func (s Server) DeletePost(ctx context.Context, request *protobuf.DeletePostRequest) (*protobuf.DeletePostResponse, error) {
	post, err := s.repo.Get(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post %d not found", request.GetId())
	}
	ID := ctx.Value("ID").(uint64)
	if post.UserID != ID {
		return nil, status.Errorf(codes.Unauthenticated, "user %d is not the owner of post %d", ID, post.ID)
	}

	err = s.repo.Delete(post.ID)
	if err != nil {
		return nil, err
	}

	resp := &protobuf.DeletePostResponse{
		Success: true,
	}

	return resp, nil
}

func (s Server) ListPosts(ctx context.Context, req *protobuf.ListPostsRequest) (*protobuf.ListPostsResponse, error) {
	list, err := s.repo.List(int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, err
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
		return nil, status.Error(codes.NotFound, "user not found")
	}
	users := userResp.GetUsers()
	var posts []*protobuf.ResponsePost
	for _, post := range list {
		user := protobuf.User{}
		for _, user := range users {
			if post.UserID == user.Id {
				user = user
			}
		}
		posts = append(posts, entityToProtobuf(post, &user))
	}

	resp := &protobuf.ListPostsResponse{
		Posts: posts,
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
