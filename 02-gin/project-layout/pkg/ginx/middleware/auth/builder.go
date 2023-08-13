package authtoken

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"project-layout/pkg/ginx"
	"project-layout/pkg/jwt"
)

// Builder 鉴权，验证用户token是否有效
type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() gin.HandlerFunc {
	return ginx.Handle(func(c *ginx.Context) {
		tokenString, err := getJwtFromHeader(c)
		if err != nil {
			c.JSONE(http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}
		// 验证token是否正确
		claims, err := jwt.ParseToken(tokenString, "secret")
		if err != nil {
			c.JSONE(http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	})
}

func getJwtFromHeader(c *ginx.Context) (string, error) {
	// 请求头的形式为 Authorization: Bearer token
	tokenString := c.Request.Header.Get("Authorization")
	if len(tokenString) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	strs := strings.SplitN(tokenString, " ", 2)
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", fmt.Errorf("token 不符合规则")
	}
	return strs[1], nil
}
