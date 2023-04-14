//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/auth/v1"
	"github.com/jxlwqq/blog-microservices/internal/auth"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
)

func InitServer(logger log.Logger, conf *config.Config) (v1.AuthServiceServer, error) {
	wire.Build(
		jwt.NewManager,
		auth.NewServer,
	)
	return auth.Server{}, nil
}
