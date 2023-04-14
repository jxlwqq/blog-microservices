package comment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	c1 := &Comment{
		UUID:    uuid.NewString(),
		UserID:  1,
		PostID:  1,
		Content: "Hello World",
	}

	c2 := &Comment{
		UUID:    uuid.NewString(),
		UserID:  2,
		PostID:  1,
		Content: "Hello Go",
	}

	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewCommentDB(conf, logger)
	require.NoError(t, err)
	repo := NewRepository(logger, db)

	// Test Create
	err = repo.Create(context.Background(), c1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), c2)
	require.NoError(t, err)

	// Test Get
	comment, err := repo.Get(context.Background(), c1.ID)
	require.NoError(t, err)
	require.Equal(t, c1.UUID, comment.UUID)

	// Test GetByUUID
	comment, err = repo.GetByUUID(context.Background(), c1.UUID)
	require.NoError(t, err)
	require.Equal(t, c1.UUID, comment.UUID)

	// Test ListByPostID
	comments, err := repo.ListByPostID(context.Background(), uint64(1), 0, 10)
	require.NoError(t, err)
	require.Equal(t, 2, len(comments))
	require.Equal(t, c1.UUID, comments[0].UUID)
	require.Equal(t, c2.UUID, comments[1].UUID)

	// Test CountByPostID
	count, err := repo.CountByPostID(context.Background(), uint64(1))
	require.NoError(t, err)
	require.Equal(t, uint64(2), count)

	// Test Update
	c1.Content = "Hello World Updated"
	err = repo.Update(context.Background(), c1)
	require.NoError(t, err)

	// Test Delete
	err = repo.Delete(context.Background(), c1.ID)
	require.NoError(t, err)

	// Test DeleteByUUID
	err = repo.DeleteByUUID(context.Background(), c2.UUID)
	require.NoError(t, err)

	// Test GetWithUnscoped
	comment, err = repo.GetWithUnscoped(context.Background(), c1.ID)
	require.NoError(t, err)
	require.Equal(t, c1.UUID, comment.UUID)
}
