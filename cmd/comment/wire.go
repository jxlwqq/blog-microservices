//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jxlwqq/blog-microservices/api/protobuf"
	"github.com/jxlwqq/blog-microservices/internal/comment"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/post"
	"github.com/jxlwqq/blog-microservices/internal/user"
)

func InitServer(logger *log.Logger, conf *config.Config) (protobuf.CommentServiceServer, error) {
	wire.Build(
		dbcontext.NewCommentDB,
		comment.NewRepository,
		user.NewClient,
		post.NewClient,
		comment.NewServer,
	)
	return &comment.Server{}, nil
}
