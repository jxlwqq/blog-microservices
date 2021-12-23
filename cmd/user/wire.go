//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jxlwqq/blog-microservices/api/protobuf/user/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/user"
)

func InitServer(logger *log.Logger, conf *config.Config) (v1.UserServiceServer, error) {
	wire.Build(
		dbcontext.NewUserDB,
		user.NewRepository,
		user.NewServer,
	)
	return &user.Server{}, nil
}
