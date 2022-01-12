package blog

import (
	"context"
	"testing"

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

func TestServer_SignUp(t *testing.T) {
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
