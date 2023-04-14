package dbcontext

import (
	"testing"

	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)

	logger := log.New()
	userDB, err := NewUserDB(conf, logger)
	require.NoError(t, err)
	require.NotNil(t, userDB)

	postDB, err := NewPostDB(conf, logger)
	require.NoError(t, err)
	require.NotNil(t, postDB)

	commentDB, err := NewCommentDB(conf, logger)
	require.NoError(t, err)
	require.NotNil(t, commentDB)

}
