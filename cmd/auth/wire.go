//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/pkg/jwt"
	"github.com/stonecutter/blog-microservices/internal/pkg/log"
	"github.com/stonecutter/blog-microservices/internal/user"
)

func InitServer(logger *log.Logger, conf *config.Config) (protobuf.AuthServiceServer, error) {
	wire.Build(
		user.NewClient,
		jwt.NewJWTManager,
		auth.NewServer,
	)
	return auth.Server{}, nil
}
