package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/user/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	u := &v1.User{
		Uuid:     uuid.NewString(),
		Username: "test1",
		Email:    "test1@test.com",
		Password: "test1",
	}

	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewUserDB(conf)
	require.NoError(t, err)
	repo := NewRepository(logger, db)
	require.NotNil(t, repo)
	s := NewServer(logger, repo)
	require.NotNil(t, s)

	createUserReq := &v1.CreateUserRequest{User: u}
	createUserResp, err := s.CreateUser(context.Background(), createUserReq)
	require.NoError(t, err)
	require.NotNil(t, createUserResp.GetUser().GetId())

	getUserByEmailReq := &v1.GetUserByEmailRequest{Email: u.GetEmail(), Password: u.GetPassword()}
	getUserByEmailResp, err := s.GetUserByEmail(context.Background(), getUserByEmailReq)
	if err != nil {
		return
	}
	require.NoError(t, err)
	require.EqualValues(t, getUserByEmailResp.GetUser().GetEmail(), u.GetEmail())
	require.NotNil(t, getUserByEmailResp.GetUser().GetId())

	getUserByUsernameReq := &v1.GetUserByUsernameRequest{Username: u.GetUsername(), Password: u.GetPassword()}
	getUserByUsernameResp, err := s.GetUserByUsername(context.Background(), getUserByUsernameReq)
	require.NoError(t, err)
	require.EqualValues(t, getUserByUsernameResp.GetUser().GetUsername(), u.GetUsername())
	require.NotNil(t, getUserByUsernameResp.GetUser().GetId())

	deleteUserReq := &v1.DeleteUserRequest{Id: createUserReq.GetUser().GetId()}
	deleteUserResp, err := s.DeleteUser(context.Background(), deleteUserReq)
	require.NoError(t, err)
	require.True(t, deleteUserResp.GetSuccess())
}
