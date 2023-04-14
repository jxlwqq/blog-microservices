//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/comment/v1"
	"github.com/jxlwqq/blog-microservices/internal/comment"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
)

func InitServer(logger log.Logger, conf *config.Config) (v1.CommentServiceServer, error) {
	wire.Build(
		dbcontext.NewCommentDB,
		comment.NewRepository,
		comment.NewServer,
	)
	return &comment.Server{}, nil
}
