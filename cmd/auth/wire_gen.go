// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/jxlwqq/blog-microservices/api/protobuf/auth/v1"
	"github.com/jxlwqq/blog-microservices/internal/auth"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
)

// Injectors from wire.go:

func InitServer(logger *log.Logger, conf *config.Config) (v1.AuthServiceServer, error) {
	manager := jwt.NewManager(logger, conf)
	authServiceServer := auth.NewServer(logger, manager)
	return authServiceServer, nil
}
