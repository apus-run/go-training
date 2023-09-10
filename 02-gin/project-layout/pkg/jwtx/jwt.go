package jwtx

import (
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken 生成jwt token
func GenerateToken(options ...Option) (string, error) {
	opts := Apply(options...)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Uid:              opts.uid,
		Ssid:             opts.ssid,
		UserAgent:        opts.userAgent,
		RegisteredClaims: opts.RegisteredClaims,
	})
	tokenString, err := token.SignedString([]byte(opts.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析jwt token
func ParseToken(tokenString, secretKey string) (*CustomClaims, *jwt.Token, error) {
	cc := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, cc, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return cc, nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, token, err
	}
	return cc, nil, err
}
