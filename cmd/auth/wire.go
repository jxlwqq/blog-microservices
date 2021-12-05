//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/auth"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/user"
)

func InitServer(conf *config.Config) (protobuf.AuthServiceServer, error) {
	wire.Build(
		user.NewClient,
		auth.NewJWTManager,
		auth.NewServer,
	)
	return auth.Server{}, nil
}
