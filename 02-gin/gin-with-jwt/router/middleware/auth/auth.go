package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
				"msg":  "没有登录",
			})
			ctx.Abort()
			return
		}
	}
}
