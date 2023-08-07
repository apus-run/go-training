package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"time"
)

// CustomClaims 在标准声明中加入用户id
type CustomClaims struct {
	UserId uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成jwt token
func GenerateToken(options ...Option) (string, error) {
	opts := Apply(options...)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserId: opts.userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(opts.expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "server",
			Subject:   "",
			ID:        "",
			Audience:  []string{},
		},
	})
	tokenString, err := token.SignedString([]byte(opts.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析jwt token
func ParseToken(tokenString, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, err
	} else {
		return nil, err
	}
}
