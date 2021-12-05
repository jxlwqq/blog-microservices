package dbcontext

import (
	"github.com/stonecutter/blog-microservices/internal/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func NewUserDB(conf *config.Config) (*DB, error) {
	return NewDB(conf.User.DB.DSN)
}

func NewPostDB(conf *config.Config) (*DB, error) {
	return NewDB(conf.Post.DB.DSN)
}

func NewCommentDB(conf *config.Config) (*DB, error) {
	return NewDB(conf.Comment.DB.DSN)
}
