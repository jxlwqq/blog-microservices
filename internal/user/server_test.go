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
	db, err := dbcontext.NewUserDB(conf, logger)
	require.NoError(t, err)
	repo := NewRepository(logger, db)
	require.NotNil(t, repo)
	s := NewServer(logger, repo)
	require.NotNil(t, s)

	// Test Create
	createUserReq := &v1.CreateUserRequest{User: u}
	createUserResp, err := s.CreateUser(context.Background(), createUserReq)
	require.NoError(t, err)
	require.NotNil(t, createUserResp.GetUser().GetId())

	// Test GetByEmail
	getUserByEmailReq := &v1.GetUserByEmailRequest{Email: u.GetEmail(), Password: u.GetPassword()}
	getUserByEmailResp, err := s.GetUserByEmail(context.Background(), getUserByEmailReq)
	if err != nil {
		return
	}
	require.NoError(t, err)
	require.EqualValues(t, getUserByEmailResp.GetUser().GetEmail(), u.GetEmail())
	require.NotNil(t, getUserByEmailResp.GetUser().GetId())

	// Test GetByUsername
	getUserByUsernameReq := &v1.GetUserByUsernameRequest{Username: u.GetUsername(), Password: u.GetPassword()}
	getUserByUsernameResp, err := s.GetUserByUsername(context.Background(), getUserByUsernameReq)
	require.NoError(t, err)
	require.EqualValues(t, getUserByUsernameResp.GetUser().GetUsername(), u.GetUsername())
	require.NotNil(t, getUserByUsernameResp.GetUser().GetId())

	// Test GetByID
	getUserReq := &v1.GetUserRequest{Id: createUserResp.GetUser().GetId()}
	getUserResp, err := s.GetUser(context.Background(), getUserReq)
	require.NoError(t, err)
	require.EqualValues(t, getUserResp.GetUser().GetId(), createUserResp.GetUser().GetId())

	// Test ListUsersByIDs
	listUsersByIDsReq := &v1.ListUsersByIDsRequest{Ids: []uint64{createUserResp.GetUser().GetId()}}
	listUsersByIDsResp, err := s.ListUsersByIDs(context.Background(), listUsersByIDsReq)
	require.NoError(t, err)
	require.Equal(t, 1, len(listUsersByIDsResp.GetUsers()))
	require.EqualValues(t, listUsersByIDsResp.GetUsers()[0].GetId(), createUserResp.GetUser().GetId())

	// Test UpdateEmail
	updateUserReq := &v1.UpdateUserRequest{User: &v1.User{Id: createUserResp.GetUser().GetId(), Email: "test123@test.com"}}
	updateUserResp, err := s.UpdateUser(context.Background(), updateUserReq)
	require.NoError(t, err)
	require.True(t, updateUserResp.GetSuccess())

	// Test UpdatePassword
	updateUserReq = &v1.UpdateUserRequest{User: &v1.User{Id: createUserResp.GetUser().GetId(), Password: "test123"}}
	updateUserResp, err = s.UpdateUser(context.Background(), updateUserReq)
	require.NoError(t, err)
	require.True(t, updateUserResp.GetSuccess())

	// Test UpdateUsername
	updateUserReq = &v1.UpdateUserRequest{User: &v1.User{Id: createUserResp.GetUser().GetId(), Username: "test123"}}
	updateUserResp, err = s.UpdateUser(context.Background(), updateUserReq)
	require.NoError(t, err)
	require.True(t, updateUserResp.GetSuccess())

	// Test UpdateAvatar
	updateUserReq = &v1.UpdateUserRequest{User: &v1.User{Id: createUserResp.GetUser().GetId(), Avatar: "test123.png"}}
	updateUserResp, err = s.UpdateUser(context.Background(), updateUserReq)
	require.NoError(t, err)
	require.True(t, updateUserResp.GetSuccess())

	// Test Delete
	deleteUserReq := &v1.DeleteUserRequest{Id: createUserReq.GetUser().GetId()}
	deleteUserResp, err := s.DeleteUser(context.Background(), deleteUserReq)
	require.NoError(t, err)
	require.True(t, deleteUserResp.GetSuccess())
}

func TestGenerateFromPassword(t *testing.T) {
	generated, err := generateFromPassword("123")
	require.NoError(t, err)
	require.NotEmpty(t, generated)
}

func TestIsCorrectPassword(t *testing.T) {
	generated, err := generateFromPassword("123")
	require.NoError(t, err)
	require.NotEmpty(t, generated)
	require.True(t, isCorrectPassword(generated, "123"))
}

func TestEntityToProtobuf(t *testing.T) {
	u := entityToProtobuf(&User{})
	require.NotNil(t, u)
}
