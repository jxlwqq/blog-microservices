//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jxlwqq/blog-microservices/api/protobuf"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/post"
	"github.com/jxlwqq/blog-microservices/internal/user"
)

func InitServer(logger *log.Logger, conf *config.Config) (protobuf.PostServiceServer, error) {
	wire.Build(
		dbcontext.NewPostDB,
		post.NewRepository,
		user.NewClient,
		post.NewServer,
	)
	return &post.Server{}, nil
}
