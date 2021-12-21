package user

import (
	"context"
	"github.com/jxlwqq/blog-microservices/api/protobuf"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
	"testing"
)

var u = &protobuf.User{
	Username: "test1",
	Email:    "test1@test.com",
	Password: "test1",
}

func newServer(t *testing.T) protobuf.UserServiceServer {
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
	return s
}

func TestServer_CreateUser(t *testing.T) {
	s := newServer(t)
	req := &protobuf.CreateUserRequest{User: u}
	resp, err := s.CreateUser(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp.GetUser().GetId())
}

func TestServer_GetUserByEmail(t *testing.T) {
	req := &protobuf.GetUserByEmailRequest{Email: u.GetEmail(), Password: u.GetPassword()}
	s := newServer(t)
	resp, err := s.GetUserByEmail(context.Background(), req)
	require.NoError(t, err)
	require.EqualValues(t, resp.GetUser().GetEmail(), u.GetEmail())
	require.NotNil(t, resp.GetUser().GetId())
}

func TestServer_GetUserByUsername(t *testing.T) {
	req := &protobuf.GetUserByUsernameRequest{Username: u.GetUsername(), Password: u.GetPassword()}
	s := newServer(t)
	resp, err := s.GetUserByUsername(context.Background(), req)
	require.NoError(t, err)
	require.EqualValues(t, resp.GetUser().GetUsername(), u.GetUsername())
	require.NotNil(t, resp.GetUser().GetId())
}

func TestServer_DeleteUser(t *testing.T) {
	s := newServer(t)
	req := &protobuf.DeleteUserRequest{Id: u.Id}
	resp, err := s.DeleteUser(context.Background(), req)
	require.NoError(t, err)
	require.True(t, resp.GetSuccess())
}
