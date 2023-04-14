//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/blog/v1"
	"github.com/jxlwqq/blog-microservices/internal/auth"
	"github.com/jxlwqq/blog-microservices/internal/blog"
	"github.com/jxlwqq/blog-microservices/internal/comment"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/post"
	"github.com/jxlwqq/blog-microservices/internal/user"
)

func InitServer(logger log.Logger, conf *config.Config) (v1.BlogServiceServer, error) {
	wire.Build(
		user.NewClient,
		auth.NewClient,
		post.NewClient,
		comment.NewClient,
		blog.NewServer,
	)

	return &blog.Server{}, nil
}
