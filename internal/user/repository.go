package user

import (
	"github.com/stonecutter/blog-microservices/internal/pkg/dbcontext"
	"github.com/stonecutter/blog-microservices/internal/pkg/log"
)

func NewRepository(logger *log.Logger, db *dbcontext.DB) Repository {
	return &repository{
		logger: logger,
		db:     db,
	}
}

type Repository interface {
	GetListByIDs(ids []uint64) ([]*User, error)
	Get(id uint64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint64) error
}

type repository struct {
	logger *log.Logger
	db     *dbcontext.DB
}

func (r repository) GetListByIDs(ids []uint64) ([]*User, error) {
	users := []*User{}
	err := r.db.Where("id IN (?)", ids).Find(&users).Error
	return users, err
}

func (r repository) Get(id uint64) (*User, error) {
	user := &User{}
	err := r.db.First(user, id).Error
	return user, err
}

func (r repository) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.Where("email = ?", email).First(user).Error
	return user, err
}

func (r repository) GetByUsername(username string) (*User, error) {
	user := &User{}
	err := r.db.Where("username = ?", username).First(user).Error
	return user, err
}

func (r repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r repository) Delete(id uint64) error {
	return r.db.Delete(&User{}, id).Error
}
