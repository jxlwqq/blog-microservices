//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jxlwqq/blog-microservices/api/protobuf/post/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/post"
)

func InitServer(logger *log.Logger, conf *config.Config) (v1.PostServiceServer, error) {
	wire.Build(
		dbcontext.NewPostDB,
		post.NewRepository,
		post.NewServer,
	)
	return &post.Server{}, nil
}
