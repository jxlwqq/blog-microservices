package user

import (
	"context"
	"fmt"

	"github.com/stonecutter/blog-microservices/api/protobuf"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServer(repo Repository) protobuf.UserServiceServer {
	return &Server{repo: repo}
}

type Server struct {
	protobuf.UnimplementedUserServiceServer
	repo Repository
}

func (s Server) GetUserListByIDs(ctx context.Context, req *protobuf.GetUserListByIDsRequest) (*protobuf.GetUserListByIDsResponse, error) {
	ids := req.GetIds()
	users, err := s.repo.GetListByIDs(ids)
	if err != nil {
		return nil, err
	}
	protoUsers := make([]*protobuf.User, len(users))
	for i, user := range users {
		protoUsers[i] = entityToProtobuf(user)
	}
	resp := &protobuf.GetUserListByIDsResponse{
		Users: protoUsers,
	}
	return resp, nil
}

func (s Server) GetUser(ctx context.Context, req *protobuf.GetUserRequest) (*protobuf.GetUserResponse, error) {
	id := req.GetId()
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	resp := &protobuf.GetUserResponse{
		User: entityToProtobuf(user),
	}

	return resp, nil
}

func (s Server) GetUserByEmail(ctx context.Context, req *protobuf.GetUserByEmailRequest) (*protobuf.GetUserResponse, error) {
	email := req.GetEmail()
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	ok := isCorrectPassword(user.Password, req.GetPassword())
	if !ok {
		return nil, status.Errorf(codes.Internal, "incorrect password")
	}
	resp := &protobuf.GetUserResponse{
		User: entityToProtobuf(user),
	}

	return resp, nil
}

func (s Server) GetUserByUsername(ctx context.Context, req *protobuf.GetUserByUsernameRequest) (*protobuf.GetUserResponse, error) {
	username := req.GetUsername()
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	ok := isCorrectPassword(user.Password, req.GetPassword())
	if !ok {
		return nil, status.Errorf(codes.Internal, "incorrect password")
	}
	resp := &protobuf.GetUserResponse{
		User: entityToProtobuf(user),
	}

	return resp, nil
}

func (s Server) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.CreateUserResponse, error) {
	fmt.Println(req.GetUser().GetPassword())
	password, err := generateFromPassword(req.GetUser().GetPassword())
	fmt.Println(password)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username: req.GetUser().GetUsername(),
		Email:    req.GetUser().GetEmail(),
		Avatar:   req.GetUser().GetAvatar(),
		Password: password,
	}
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	resp := &protobuf.CreateUserResponse{
		User: entityToProtobuf(user),
	}

	return resp, nil
}

func (s Server) UpdateUser(ctx context.Context, req *protobuf.UpdateUserRequest) (*protobuf.UpdateUserResponse, error) {

	user := &User{ID: req.GetUser().GetId()}
	if req.GetUser().GetUsername() != "" {
		user.Username = req.GetUser().GetUsername()
	}
	if req.GetUser().GetPassword() != "" {
		password, err := generateFromPassword(req.GetUser().GetPassword())
		if err != nil {
			return nil, err
		}
		user.Password = password
	}
	if req.GetUser().Email != "" {
		user.Email = req.GetUser().Email
	}
	if req.GetUser().Avatar != "" {
		user.Avatar = req.GetUser().Avatar
	}
	err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	resp := &protobuf.UpdateUserResponse{
		Success: true,
	}

	return resp, nil
}

func generateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func isCorrectPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func entityToProtobuf(user *User) *protobuf.User {
	return &protobuf.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
