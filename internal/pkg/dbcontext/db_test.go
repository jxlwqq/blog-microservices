package dbcontext

import (
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)

	userDB, err := NewUserDB(conf)
	require.NoError(t, err)
	require.NotNil(t, userDB)

	postDB, err := NewPostDB(conf)
	require.NoError(t, err)
	require.NotNil(t, postDB)

	commentDB, err := NewCommentDB(conf)
	require.NoError(t, err)
	require.NotNil(t, commentDB)

}
