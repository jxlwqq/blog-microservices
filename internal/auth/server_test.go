package auth

import (
	"context"
	"testing"

	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/auth/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	jwtManager := jwt.NewManager(logger, conf)
	server := NewServer(logger, jwtManager)

	genResp, err := server.GenerateToken(context.Background(), &v1.GenerateTokenRequest{
		UserId: 1,
	})
	if err != nil {
		return
	}
	require.NoError(t, err)
	require.NotEmpty(t, genResp.GetToken())

	validateResp, err := server.ValidateToken(context.Background(), &v1.ValidateTokenRequest{
		Token: genResp.GetToken(),
	})
	require.NoError(t, err)
	require.True(t, validateResp.GetValid())

	refreshResp, err := server.RefreshToken(context.Background(), &v1.RefreshTokenRequest{
		Token: genResp.GetToken(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, refreshResp.GetToken())

	refreshResp, err = server.RefreshToken(context.Background(), &v1.RefreshTokenRequest{
		Token: "a.b.c",
	})
	require.Error(t, err)
	require.Nil(t, refreshResp)
}
