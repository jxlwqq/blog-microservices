package blog

import (
	"context"
	"testing"

	commentv1 "github.com/jxlwqq/blog-microservices/api/protobuf/comment/v1"

	postv1 "github.com/jxlwqq/blog-microservices/api/protobuf/post/v1"

	authv1 "github.com/jxlwqq/blog-microservices/api/protobuf/auth/v1"

	"github.com/google/uuid"
	userv1 "github.com/jxlwqq/blog-microservices/api/protobuf/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/blog/v1"

	"github.com/golang/mock/gomock"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/mock"
	"github.com/stretchr/testify/require"
)

func makeMock(t *testing.T) (v1.BlogServiceServer, *mock.MockUserServiceClient, *mock.MockPostServiceClient, *mock.MockCommentServiceClient, *mock.MockAuthServiceClient) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserClient := mock.NewMockUserServiceClient(ctl)
	mockPostClient := mock.NewMockPostServiceClient(ctl)
	mockCommentClient := mock.NewMockCommentServiceClient(ctl)
	mockAuthClient := mock.NewMockAuthServiceClient(ctl)
	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	s := NewServer(logger, conf, mockUserClient, mockPostClient, mockCommentClient, mockAuthClient)
	require.NotNil(t, s)

	return s, mockUserClient, mockPostClient, mockCommentClient, mockAuthClient
}

func TestServer_SignUp(t *testing.T) {
	s, mockUserClient, _, _, mockAuthClient := makeMock(t)
	gomock.InOrder(
		mockUserClient.EXPECT().GetUserByUsername(context.Background(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "")),
		mockUserClient.EXPECT().GetUserByEmail(context.Background(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "")),
		mockUserClient.EXPECT().CreateUser(context.Background(), gomock.Any()).Return(&userv1.CreateUserResponse{User: &userv1.User{
			Id:       1,
			Uuid:     uuid.NewString(),
			Username: "test",
			Email:    "test@test.com",
		}}, nil),
		mockAuthClient.EXPECT().GenerateToken(context.Background(), gomock.Any()).Return(&authv1.GenerateTokenResponse{Token: "token"}, nil),
	)

	resp, err := s.SignUp(context.Background(), &v1.SignUpRequest{
		Username: "test",
		Email:    "test@test.com",
		Password: "pass",
	})
	require.NoError(t, err)

	require.Equal(t, "token", resp.GetToken())
}

func TestServer_SignIn_WithUsername(t *testing.T) {
	s, mockUserClient, _, _, mockAuthClient := makeMock(t)
	gomock.InOrder(
		mockUserClient.EXPECT().GetUserByUsername(context.Background(), gomock.Any()).Return(&userv1.GetUserResponse{User: &userv1.User{
			Id:       1,
			Uuid:     uuid.NewString(),
			Username: "test",
			Email:    "test@test.com",
		}}, nil),
		mockAuthClient.EXPECT().GenerateToken(context.Background(), gomock.Any()).Return(&authv1.GenerateTokenResponse{Token: "token"}, nil),
	)

	resp, err := s.SignIn(context.Background(), &v1.SignInRequest{
		Request: &v1.SignInRequest_Username{
			Username: "test",
		},
		Password: "pass",
	})
	require.NoError(t, err)

	require.Equal(t, "token", resp.GetToken())
}

func TestServer_SignIn_WithEmail(t *testing.T) {
	s, mockUserClient, _, _, mockAuthClient := makeMock(t)
	gomock.InOrder(
		mockUserClient.EXPECT().GetUserByEmail(context.Background(), gomock.Any()).Return(&userv1.GetUserResponse{User: &userv1.User{
			Id:       1,
			Uuid:     uuid.NewString(),
			Username: "test",
			Email:    "test@test.com",
		}}, nil),
		mockAuthClient.EXPECT().GenerateToken(context.Background(), gomock.Any()).Return(&authv1.GenerateTokenResponse{Token: "token"}, nil),
	)

	resp, err := s.SignIn(context.Background(), &v1.SignInRequest{
		Request: &v1.SignInRequest_Email{
			Email: "test@test.com",
		},
		Password: "pass",
	})
	require.NoError(t, err)

	require.Equal(t, "token", resp.GetToken())
}

func TestServer_CreatePost(t *testing.T) {
	s, mockUserClient, mockPostClient, _, _ := makeMock(t)

	ctx := context.WithValue(context.Background(), "ID", uint64(1))

	gomock.InOrder(
		mockUserClient.EXPECT().GetUser(ctx, gomock.Any()).Return(&userv1.GetUserResponse{User: &userv1.User{
			Id:       1,
			Uuid:     uuid.NewString(),
			Username: "test",
			Email:    "test@test.com",
		}}, nil),
		mockPostClient.EXPECT().CreatePost(ctx, gomock.Any()).Return(&postv1.CreatePostResponse{Post: &postv1.Post{
			Id:      uint64(1),
			Uuid:    uuid.NewString(),
			Title:   "test",
			Content: "test",
		}}, nil),
	)

	createResp, err := s.CreatePost(ctx, &v1.CreatePostRequest{
		Post: &v1.Post{
			Title:   "test",
			Content: "test",
		},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), createResp.GetPost().GetId())
}

func TestServer_GetPost(t *testing.T) {
	s, mockUserClient, mockPostClient, _, _ := makeMock(t)

	gomock.InOrder(
		mockPostClient.EXPECT().GetPost(context.Background(), gomock.Any()).Return(&postv1.GetPostResponse{Post: &postv1.Post{
			Id:      uint64(1),
			Uuid:    uuid.NewString(),
			Title:   "test",
			Content: "test",
		}}, nil),
		mockUserClient.EXPECT().GetUser(context.Background(), gomock.Any()).Return(&userv1.GetUserResponse{User: &userv1.User{
			Id:       1,
			Uuid:     uuid.NewString(),
			Username: "test",
			Email:    "test@test.com",
		}}, nil),
	)

	getResp, err := s.GetPost(context.Background(), &v1.GetPostRequest{
		Id: uint64(1),
	})
	if err != nil {
		return
	}

	require.Equal(t, uint64(1), getResp.GetPost().GetId())
}

func TestServer_UpdatePost(t *testing.T) {
	s, mockUserClient, mockPostClient, _, _ := makeMock(t)
	ctx := context.WithValue(context.Background(), "ID", uint64(1))
	gomock.InOrder(
		mockUserClient.EXPECT().GetUser(ctx, gomock.Any()).Return(&userv1.GetUserResponse{User: &userv1.User{
			Id:       uint64(1),
			Uuid:     uuid.NewString(),
			Username: "test",
			Email:    "test@test.com",
		}}, nil),
		mockPostClient.EXPECT().GetPost(ctx, gomock.Any()).Return(&postv1.GetPostResponse{Post: &postv1.Post{
			Id:      uint64(1),
			Uuid:    uuid.NewString(),
			UserId:  uint64(1),
			Title:   "test",
			Content: "test",
		}}, nil),
		mockPostClient.EXPECT().UpdatePost(ctx, gomock.Any()).Return(&postv1.UpdatePostResponse{Success: true}, nil),
	)
	updateResp, err := s.UpdatePost(ctx, &v1.UpdatePostRequest{
		Post: &v1.Post{
			Id:      uint64(1),
			Title:   "test2",
			Content: "test2",
		},
	})
	require.NoError(t, err)
	require.True(t, updateResp.GetSuccess())
}

func TestServer_ListPosts(t *testing.T) {
	s, mockUserClient, mockPostClient, _, _ := makeMock(t)

	gomock.InOrder(
		mockPostClient.EXPECT().ListPosts(context.Background(), gomock.Any()).Return(&postv1.ListPostsResponse{
			Posts: []*postv1.Post{
				{
					Id:      uint64(1),
					Uuid:    uuid.NewString(),
					UserId:  uint64(1),
					Title:   "test",
					Content: "test",
				},
				{
					Id:      uint64(2),
					Uuid:    uuid.NewString(),
					UserId:  uint64(2),
					Title:   "test2",
					Content: "test2",
				},
			},
			Count: 2,
		}, nil),
		mockUserClient.EXPECT().ListUsersByIDs(context.Background(), gomock.Any()).Return(&userv1.ListUsersByIDsResponse{
			Users: []*userv1.User{
				{
					Id:       uint64(1),
					Uuid:     uuid.NewString(),
					Username: "test",
					Email:    "test@test.com",
				},
				{
					Id:       uint64(2),
					Uuid:     uuid.NewString(),
					Username: "test2",
					Email:    "test2@test.com",
				},
			},
		}, nil),
	)

	listResp, err := s.ListPosts(context.Background(), &v1.ListPostsRequest{
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(2), listResp.GetTotal())
	require.Equal(t, uint64(1), listResp.GetPosts()[0].GetUser().GetId())
	require.Equal(t, uint64(2), listResp.GetPosts()[1].GetUser().GetId())
}

func TestServer_UpdateComment(t *testing.T) {
	s, mockUserClient, _, mockCommentClient, _ := makeMock(t)

	ctx := context.WithValue(context.Background(), "ID", uint64(1))

	gomock.InOrder(
		mockUserClient.EXPECT().GetUser(ctx, gomock.Any()).Return(&userv1.GetUserResponse{User: &userv1.User{
			Id:       uint64(1),
			Uuid:     uuid.NewString(),
			Username: "test",
			Email:    "test@test.com",
		}}, nil),
		mockCommentClient.EXPECT().GetComment(ctx, gomock.Any()).Return(&commentv1.GetCommentResponse{Comment: &commentv1.Comment{
			Id:      uint64(1),
			Uuid:    uuid.NewString(),
			UserId:  uint64(1),
			Content: "Hello World",
		}}, nil),
		mockCommentClient.EXPECT().UpdateComment(ctx, gomock.Any()).Return(&commentv1.UpdateCommentResponse{Success: true}, nil),
	)

	updateResp, err := s.UpdateComment(ctx, &v1.UpdateCommentRequest{Comment: &v1.Comment{
		Id:      1,
		Content: "Hello Go",
	}})

	require.NoError(t, err)
	require.True(t, updateResp.GetSuccess())
}

func TestServer_ListCommentsByPostID(t *testing.T) {
	s, mockUserClient, _, mockCommentClient, _ := makeMock(t)

	gomock.InOrder(
		mockCommentClient.EXPECT().ListCommentsByPostID(context.Background(), gomock.Any()).Return(&commentv1.ListCommentsByPostIDResponse{
			Comments: []*commentv1.Comment{
				{
					Id:      uint64(1),
					Uuid:    uuid.NewString(),
					UserId:  uint64(1),
					Content: "Hello World",
				},
				{
					Id:      uint64(2),
					Uuid:    uuid.NewString(),
					UserId:  uint64(2),
					Content: "Hello Go",
				},
			},
			Total: 2,
		}, nil),
		mockUserClient.EXPECT().ListUsersByIDs(context.Background(), gomock.Any()).Return(&userv1.ListUsersByIDsResponse{
			Users: []*userv1.User{
				{
					Id:       uint64(1),
					Uuid:     uuid.NewString(),
					Username: "test",
					Email:    "test@test.com",
				},
				{
					Id:       uint64(2),
					Uuid:     uuid.NewString(),
					Username: "test2",
					Email:    "test2@test.com",
				},
			},
		}, nil),
	)

	listResp, err := s.ListCommentsByPostID(context.Background(), &v1.ListCommentsByPostIDRequest{
		PostId: 1,
		Offset: 0,
		Limit:  10,
	})

	require.NoError(t, err)
	require.Equal(t, uint64(2), listResp.GetTotal())
	require.Equal(t, uint64(1), listResp.GetComments()[0].GetUser().GetId())
	require.Equal(t, uint64(2), listResp.GetComments()[1].GetUser().GetId())

}
