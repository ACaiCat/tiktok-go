package jwt

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID    int64 `json:"user_id"`
	TokenType int8  `json:"token_type"`
	jwt.RegisteredClaims
}
