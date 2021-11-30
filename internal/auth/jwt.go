package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

func NewJWTManager(secret string, expires time.Duration) *JWTManager {
	return &JWTManager{
		secret:  secret,
		expires: expires,
	}
}

type JWTManager struct {
	secret  string
	expires time.Duration
}

type UserClaims struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (manager *JWTManager) Generate(id uint64, username string) (string, error) {
	claims := UserClaims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.expires).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secret))
}
