package user

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
