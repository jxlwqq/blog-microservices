package post

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID            uint64         `json:"id"`
	UUID          string         `json:"uuid"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	CommentsCount uint32         `json:"comments_count"`
	UserID        uint64         `json:"user_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}
