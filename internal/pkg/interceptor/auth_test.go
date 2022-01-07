package interceptor

import (
	"context"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestAuthInterceptor(t *testing.T) {
	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	jwtManager := jwt.NewManager(logger, conf)
	methods := map[string]bool{
		"Create": true,
		"Get":    false,
	}
	i := NewAuthInterceptor(logger, jwtManager, methods)

	_, err = i.Authorize(context.Background(), "Create")
	require.Error(t, err)

	token, err := jwtManager.Generate(uint64(1))
	require.NoError(t, err)

	header := metadata.New(map[string]string{headerAuthorize: expectedScheme + " " + token})
	ctx := metadata.NewIncomingContext(context.Background(), header)
	_, err = i.Authorize(ctx, "Create")
	require.NoError(t, err)

	_, err = i.Authorize(ctx, "Get")
	require.NoError(t, err)
}
