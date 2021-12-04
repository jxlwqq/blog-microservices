package comment

import "github.com/stonecutter/blog-microservices/pkg/dbcontext"

func NewRepository(db *dbcontext.DB) Repository {
	return repository{
		db: db,
	}
}

type Repository interface {
	Create(comment *Comment) error
	Update(comment *Comment) error
	Delete(id uint64) error
	ListByPostID(postID uint64) ([]*Comment, error)
}

type repository struct {
	db *dbcontext.DB
}

func (r repository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

func (r repository) Update(comment *Comment) error {
	return r.db.Save(comment).Error
}

func (r repository) Delete(id uint64) error {
	return r.db.Delete(&Comment{ID: id}).Error
}

func (r repository) ListByPostID(postID uint64) ([]*Comment, error) {
	var comments []*Comment
	err := r.db.Where("post_id = ?", postID).Find(&comments).Error
	return comments, err
}
