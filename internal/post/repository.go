package post

import (
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
	Get(id uint64) (*Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id uint64) error
	List(offset, limit int) ([]*Post, error)
	Count() (uint64, error)
}

type repository struct {
	logger *log.Logger
	db     *dbcontext.DB
}

func (r repository) Get(id uint64) (*Post, error) {
	post := &Post{}
	err := r.db.First(post, id).Error
	return post, err
}

func (r repository) Create(post *Post) error {
	return r.db.Create(post).Error
}

func (r repository) Update(post *Post) error {
	return r.db.Save(post).Error
}

func (r repository) Delete(id uint64) error {
	return r.db.Delete(&Post{}, id).Error
}

func (r repository) List(offset, limit int) ([]*Post, error) {
	var posts []*Post
	err := r.db.Offset(offset).Limit(limit).Find(&posts).Error
	return posts, err
}

func (r repository) Count() (uint64, error) {
	var count int64
	err := r.db.Model(&Post{}).Count(&count).Error
	return uint64(count), err
}
