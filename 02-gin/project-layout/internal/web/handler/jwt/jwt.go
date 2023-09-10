package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"project-layout/pkg/ginx"
	"project-layout/pkg/jwtx"
)

type JwtHandler struct{}

func NewJwtHandler() Handler {
	return &JwtHandler{}
}

func (j *JwtHandler) SetLoginToken(ctx *ginx.Context, uid uint64) error {
	// 测试使用1分钟
	// expireAt := time.Now().Add(time.Minute)
	// 正常设置为30分钟，要将过期时间设置更长一些
	expireAt := time.Now().Add(time.Minute * 30)
	tokenStr, err := jwtx.GenerateToken(
		jwtx.WithUserAgent(ctx.GetHeader("User-Agent")),
		jwtx.WithUid(uid),
		jwtx.WithJwtRegisteredClaims(func(opts *jwtx.Options) {
			opts.ExpiresAt = jwt.NewNumericDate(expireAt)
		}),
	)

	if err != nil {
		return err
	}
	// 将token放入响应头中
	ctx.Header("x-jwt-token", tokenStr)

	return nil
}
