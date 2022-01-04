package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint64         `json:"id"`
	UUID      string         `json:"uuid"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Avatar    string         `json:"avatar"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
