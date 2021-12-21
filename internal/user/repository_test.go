package user

import (
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepository(t *testing.T) {
	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewUserDB(conf)
	require.NoError(t, err)
	repo := NewRepository(logger, db)
	require.NotNil(t, repo)

	// Test Create
	u := &User{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
		Avatar:   "https://test.com/avatar.png",
	}
	err = repo.Create(u)
	require.NoError(t, err)
	require.NotEmpty(t, u.ID)

	// Test Get
	u2, err := repo.Get(u.ID)
	require.NoError(t, err)
	require.Equal(t, u.Username, u2.Username)
	require.Equal(t, u.Email, u2.Email)

	// Test Update
	u.Avatar = "https://test.com/avatar2.png"
	err = repo.Update(u)
	require.NoError(t, err)

	// Test GetByEmail
	u3, err := repo.GetByEmail(u.Email)
	require.NoError(t, err)
	require.Equal(t, u.Email, u3.Email)

	// Test GetByUsername
	u4, err := repo.GetByUsername(u.Username)
	require.NoError(t, err)
	require.Equal(t, u.Username, u4.Username)

	// Test Delete
	err = repo.Delete(u.ID)
	require.NoError(t, err)
}
