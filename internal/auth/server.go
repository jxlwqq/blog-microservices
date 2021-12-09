package auth

import (
	"context"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/pkg/jwt"
	"github.com/stonecutter/blog-microservices/internal/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServer(logger *log.Logger, jwtManager *jwt.JWTManager, userClient protobuf.UserServiceClient) protobuf.AuthServiceServer {
	return &Server{
		logger:     logger,
		jwtManager: jwtManager,
		userClient: userClient,
	}
}

type Server struct {
	protobuf.UnimplementedAuthServiceServer
	logger     *log.Logger
	jwtManager *jwt.JWTManager
	userClient protobuf.UserServiceClient
}

func (s Server) SignIn(ctx context.Context, req *protobuf.SignInRequest) (*protobuf.SignInResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	email := req.GetEmail()
	username := req.GetUsername()
	password := req.GetPassword()
	if email == "" && username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and username cannot be empty")
	}
	var userID uint64
	var userName string
	if email != "" {
		resp, err := s.userClient.GetUserByEmail(ctx, &protobuf.GetUserByEmailRequest{
			Email:    email,
			Password: password,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
		}
		user := resp.GetUser()
		userID = user.GetId()
		userName = user.GetUsername()
	} else {
		req, err := s.userClient.GetUserByUsername(ctx, &protobuf.GetUserByUsernameRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get user by username: %v", err)
		}
		user := req.GetUser()
		userID = user.GetId()
		userName = user.GetUsername()
	}

	token, err := s.jwtManager.Generate(userID, userName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &protobuf.SignInResponse{Token: token}, nil
}

func (s Server) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error) {
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	username := req.GetUsername()
	email := req.GetEmail()
	password := req.GetPassword()

	u := &protobuf.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	resp, err := s.userClient.CreateUser(ctx, &protobuf.CreateUserRequest{
		User: u,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	userID := resp.GetUser().GetId()
	username = resp.GetUser().GetUsername()
	token, err := s.jwtManager.Generate(userID, username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	return &protobuf.SignUpResponse{Token: token}, nil

}
