package user

import (
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"github.com/stonecutter/blog-microservices/internal/pkg/dbcontext"
	"github.com/stretchr/testify/require"
	"testing"
)

func newRepository(t *testing.T) Repository {
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewUserDB(conf)
	require.NoError(t, err)
	repo := NewRepository(db)
	require.NotNil(t, repo)
	return repo
}

func TestRepository_Create(t *testing.T) {
	repo := newRepository(t)
	u := &User{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
		Avatar:   "https://test.com/avatar.png",
	}
	err := repo.Create(u)
	require.NoError(t, err)
	require.NotEmpty(t, u.ID)
}
