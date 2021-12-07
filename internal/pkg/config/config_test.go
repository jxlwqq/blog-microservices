package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetPath(t *testing.T) {
	path := GetPath()
	_, err := os.Stat(path)
	require.NoError(t, err)
}

func TestLoad(t *testing.T) {
	path := GetPath()
	conf, err := Load(path)
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
