package user

import "github.com/stonecutter/blog-microservices/pkg/dbcontext"

func NewRepository(db *dbcontext.DB) Repository {
	return &repository{
		db: db,
	}
}

type Repository interface {
	Get(id uint64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
}

type repository struct {
	db *dbcontext.DB
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
