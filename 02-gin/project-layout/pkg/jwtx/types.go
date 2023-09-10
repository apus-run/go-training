package jwtx

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtToken interface {
	GenerateToken(options ...Option) (string, error)
	ParseToken(tokenString, secretKey string) (*CustomClaims, *jwt.Token, error)
}

// SecretKey jwtx secret key
var SecretKey = "moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"

// CustomClaims 在标准声明中加入 自定义id
type CustomClaims struct {
	Uid  uint64
	Ssid string

	// UserAgent 增强安全性，防止token被盗用
	UserAgent string

	jwt.RegisteredClaims
}
