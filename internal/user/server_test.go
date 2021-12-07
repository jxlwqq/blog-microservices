package user

import (
	"context"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/pkg/dbcontext"
	"github.com/stretchr/testify/require"
	"testing"
)

var u = &protobuf.User{
	Username: "test1",
	Email:    "test1@test.com",
	Password: "test1",
}

func newServer(t *testing.T) protobuf.UserServiceServer {
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewUserDB(conf)
	require.NoError(t, err)
	repo := NewRepository(db)
	require.NotNil(t, repo)
	s := NewServer(repo)
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
