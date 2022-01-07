package user

import (
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	client, err := NewClient(logger, conf)
	require.NoError(t, err)
	require.NotNil(t, client)
}
