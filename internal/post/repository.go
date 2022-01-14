package post

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
	Get(ctx context.Context, id uint64) (*Post, error)
	GetWithUnscoped(ctx context.Context, id uint64) (*Post, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
	UpdateWithUnscoped(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id uint64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, offset, limit int) ([]*Post, error)
	Count(ctx context.Context) (uint64, error)
}

type repository struct {
	logger *log.Logger
	db     *dbcontext.DB
}

func (r repository) Get(ctx context.Context, id uint64) (*Post, error) {
	post := &Post{}
	err := r.db.First(post, id).Error
	return post, err
}

func (r repository) GetWithUnscoped(ctx context.Context, id uint64) (*Post, error) {
	post := &Post{}
	err := r.db.Unscoped().First(post, id).Error
	return post, err
}

func (r repository) Create(ctx context.Context, post *Post) error {
	return r.db.Create(post).Error
}

func (r repository) Update(ctx context.Context, post *Post) error {
	return r.db.Save(post).Error
}

func (r repository) UpdateWithUnscoped(ctx context.Context, post *Post) error {
	return r.db.Unscoped().Save(post).Error
}

func (r repository) Delete(ctx context.Context, id uint64) error {
	return r.db.Delete(&Post{}, id).Error
}

func (r repository) DeleteByUUID(ctx context.Context, uuid string) error {
	return r.db.Delete(&Post{}, "uuid = ?", uuid).Error
}

func (r repository) List(ctx context.Context, offset, limit int) ([]*Post, error) {
	var posts []*Post
	err := r.db.Offset(offset).Limit(limit).Find(&posts).Error
	return posts, err
}

func (r repository) Count(ctx context.Context) (uint64, error) {
	var count int64
	err := r.db.Model(&Post{}).Count(&count).Error
	return uint64(count), err
}
