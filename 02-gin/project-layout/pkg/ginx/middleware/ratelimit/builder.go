package ratelimit

import (
	"github.com/gin-gonic/gin"
)

type Builder struct {
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 限流逻辑

		c.Next()
	}
}
