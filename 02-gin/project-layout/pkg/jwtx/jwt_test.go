package jwtx

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	const secretKey = "moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"
	tokenStr, err := GenerateToken(
		WithUid(1),
		WithSsid("web-server"),
		WithSecretKey(secretKey),
		WithJwtRegisteredClaims(func(opts *Options) {
			//  过期时间 1 分钟
			expireAt := time.Now().Add(time.Minute * 1)
			opts.ExpiresAt = jwt.NewNumericDate(expireAt)
		}),
	)
	if err != nil {
		t.Errorf("生成 token 失败: %v", err)
		return
	}
	t.Logf("token: %v", tokenStr)
	claims, token, err := ParseToken(tokenStr, secretKey)
	if err != nil {
		t.Errorf("解析 token 失败: %v", err)
		return
	}
	t.Logf("token: %v", token)
	t.Logf("claims: %v", claims.Uid)

	// time.Sleep(time.Minute * 1)

	// 验证token是否过期
	expireTime, err := claims.GetExpirationTime()
	if err != nil {
		t.Errorf("拿不到过期时间: %v", err)
		return
	}
	now := time.Now()
	if expireTime.Before(now) {
		t.Logf("已过期: %v", err)
	}
	if claims.ExpiresAt.Time.Sub(now) < time.Second*50 {
		// 快速生成新token
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		newToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			t.Errorf("生成 新 token 失败: %v", err)
			return
		}
		t.Logf("new token: %v", newToken)
	}
}
