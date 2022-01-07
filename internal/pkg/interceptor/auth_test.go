package interceptor

import (
	"context"
	"testing"

	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/jwt"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
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
	require.NotNil(t, i)

	unary := i.Unary()
	require.NotNil(t, unary)

	_, err = i.Authorize(context.Background(), "Create")
	require.Error(t, err)

	token, err := jwtManager.Generate(uint64(1))
	require.NoError(t, err)

	header := metadata.New(map[string]string{headerAuthorize: expectedScheme + " " + token})
	ctx := metadata.NewIncomingContext(context.Background(), header)
	_, err = i.Authorize(ctx, "Create")
	require.NoError(t, err)

	header = metadata.New(map[string]string{headerAuthorize: ""})
	ctx = metadata.NewIncomingContext(context.Background(), header)
	_, err = i.Authorize(ctx, "Create")
	require.Error(t, err)

	header = metadata.New(map[string]string{headerAuthorize: expectedScheme + token})
	ctx = metadata.NewIncomingContext(context.Background(), header)
	_, err = i.Authorize(ctx, "Create")
	require.Error(t, err)

	header = metadata.New(map[string]string{headerAuthorize: "hello" + " " + token})
	ctx = metadata.NewIncomingContext(context.Background(), header)
	_, err = i.Authorize(ctx, "Create")
	require.Error(t, err)

	header = metadata.New(map[string]string{headerAuthorize: expectedScheme + " " + "a.b.c"})
	ctx = metadata.NewIncomingContext(context.Background(), header)
	_, err = i.Authorize(ctx, "Create")
	require.Error(t, err)

	_, err = i.Authorize(ctx, "Get")
	require.NoError(t, err)
}
