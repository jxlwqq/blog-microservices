// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/jxlwqq/blog-microservices/api/protobuf/user/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/jxlwqq/blog-microservices/internal/user"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

func InitServer(logger log.Logger, conf *config.Config) (v1.UserServiceServer, error) {
	db, err := dbcontext.NewUserDB(conf, logger)
	if err != nil {
		return nil, err
	}
	repository := user.NewRepository(logger, db)
	userServiceServer := user.NewServer(logger, repository)
	return userServiceServer, nil
}
