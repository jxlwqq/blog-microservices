package user

import (
	"context"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
)

func NewRepository(logger *log.Logger, db *dbcontext.DB) Repository {
	return &repository{
		logger: logger,
		db:     db,
	}
}

type Repository interface {
	ListUsersByIDs(ctx context.Context, ids []uint64) ([]*User, error)
	Get(ctx context.Context, id uint64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint64) error
}

type repository struct {
	logger *log.Logger
	db     *dbcontext.DB
}

func (r repository) ListUsersByIDs(ctx context.Context, ids []uint64) ([]*User, error) {
	users := []*User{}
	err := r.db.Where("id IN (?)", ids).Find(&users).Error
	return users, err
}

func (r repository) Get(ctx context.Context, id uint64) (*User, error) {
	user := &User{}
	err := r.db.First(user, id).Error
	return user, err
}

func (r repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.db.Where("email = ?", email).First(user).Error
	return user, err
}

func (r repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	err := r.db.Where("username = ?", username).First(user).Error
	return user, err
}

func (r repository) Create(ctx context.Context, user *User) error {
	return r.db.Create(user).Error
}

func (r repository) Update(ctx context.Context, user *User) error {
	return r.db.Save(user).Error
}

func (r repository) Delete(ctx context.Context, id uint64) error {
	return r.db.Delete(&User{}, id).Error
}
