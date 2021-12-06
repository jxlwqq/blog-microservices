package config_test

import (
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetPath(t *testing.T) {
	path := config.GetPath()
	_, err := os.Stat(path)
	require.NoError(t, err)
}

func TestLoad(t *testing.T) {
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	require.NotNil(t, conf)
	require.NotEmpty(t, conf.User.DB.DSN)
	require.NotEmpty(t, conf.Post.DB.DSN)
	require.NotEmpty(t, conf.Comment.DB.DSN)
	require.NotEmpty(t, conf.User.Server.Addr)
	require.NotEmpty(t, conf.Post.Server.Addr)
	require.NotEmpty(t, conf.Comment.Server.Addr)
	require.NotEmpty(t, conf.Auth.Server.Addr)
}
