package auth

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

var TokenKey = "moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"

// Builder 鉴权，验证用户token是否有效
type Builder struct {
	// 白名单路由地址集合, 放行
	whitePathList []string
}

func NewBuilder() *Builder {
	return &Builder{
		whitePathList: []string{},
	}
}

func (b *Builder) IgnorePaths(whitePath string) *Builder {
	b.whitePathList = append(b.whitePathList, whitePath)
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("middleware auth")

		// 白名单路由放行
		for _, path := range b.whitePathList {
			if strings.Contains(ctx.Request.URL.Path, path) {
				ctx.Next()
				return
			}
		}

		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token不能为空",
			})
			ctx.Abort()
			return
		}

		segs := strings.SplitN(tokenString, " ", 2)
		if len(segs) != 2 || segs[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token格式错误",
			})
			ctx.Abort()
			return
		}
		tokenStr := segs[1]
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(TokenKey), nil
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  err.Error(),
			})
			ctx.Abort()
			return
		}
		if token == nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token 无效",
			})
			ctx.Abort()
			return
		}

		now := time.Now()
		// 是否过期了
		if claims.ExpiresAt.Time.Before(now) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token 已过期",
			})
			ctx.Abort()
			return
		}

		// 利用 User-Agent 验证 token 是否被盗用
		if claims.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "有人在伪造身份",
			})
			ctx.Abort()
			return
		}

		// 每 10 秒刷新一次, Sub = t-d
		if claims.ExpiresAt.Time.Sub(now) < time.Second*50 {
			// 刷新
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte(TokenKey))
			if err != nil {
				log.Printf("换取 token 失败: %v", err)
				return
			}
			ctx.Header("x-jwt-token", tokenStr)
			log.Println("刷新了 token")

			return
		}

		// 将claims信息存入上下文
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

type CustomClaims struct {
	jwt.RegisteredClaims

	// 准备在 JWT 里面放一个 user id
	UserID uint64

	// 用于验证 token 是否被盗用
	UserAgent string
}
