package post

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	p1 := &Post{
		UUID:    uuid.NewString(),
		Title:   "Hello World",
		Content: "Hello World",
	}

	p2 := &Post{
		UUID:    uuid.NewString(),
		Title:   "Hello Go",
		Content: "Hello Go",
	}
	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewPostDB(conf)
	require.NoError(t, err)
	repo := NewRepository(logger, db)

	// Test Create
	err = repo.Create(context.Background(), p1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), p2)
	require.NoError(t, err)

	// Test List
	posts, err := repo.List(context.Background(), 0, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(posts))
	require.Equal(t, p1.UUID, posts[0].UUID)
	require.Equal(t, p2.UUID, posts[1].UUID)

	// Test Count
	count, err := repo.Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), count)

	// Test Get
	post, err := repo.Get(context.Background(), p1.ID)
	require.NoError(t, err)
	require.Equal(t, p1.UUID, post.UUID)

	// Test Update
	p1.Title = "Hello World2"
	err = repo.Update(context.Background(), p1)
	require.NoError(t, err)

	// Test Delete
	err = repo.Delete(context.Background(), p1.ID)
	require.NoError(t, err)

	// Test DeleteByUUID
	err = repo.DeleteByUUID(context.Background(), p2.UUID)
	require.NoError(t, err)

	// Test GetWithUnscoped
	deletedPost, err := repo.GetWithUnscoped(context.Background(), p1.ID)
	require.NoError(t, err)
	require.Equal(t, p1.ID, deletedPost.ID)

	// Test UpdateWithUnscoped
	p1.DeletedAt = gorm.DeletedAt{}
	err = repo.UpdateWithUnscoped(context.Background(), p1)
	require.NoError(t, err)
	err = repo.Delete(context.Background(), p1.ID)
	require.NoError(t, err)
}
