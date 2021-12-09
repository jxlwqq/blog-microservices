//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/stonecutter/blog-microservices/api/protobuf"
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/pkg/dbcontext"
	"github.com/stonecutter/blog-microservices/internal/pkg/log"
	"github.com/stonecutter/blog-microservices/internal/user"
)

func InitServer(logger *log.Logger, conf *config.Config) (protobuf.UserServiceServer, error) {
	wire.Build(
		dbcontext.NewUserDB,
		user.NewRepository,
		user.NewServer,
	)
	return &user.Server{}, nil
}
