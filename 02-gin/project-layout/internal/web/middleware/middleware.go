package middleware

import (
	"github.com/gin-gonic/gin"
	"project-layout/pkg/ginx"
)

// middleware 实现Router接口
// 便于服务启动时加载, middleware本质跟handler无区别
type middleware struct{}

func NewMiddleware() *middleware {
	return &middleware{}
}

// Load 注册中间件和公共路由
func (m *middleware) Load(g *gin.Engine) {
	// 注册中间件
	g.Use(ginx.NoCache())
	g.Use(ginx.CORS())
	g.Use(ginx.Secure())

	// 404
	g.NoRoute(ginx.Handle(func(c *ginx.Context) {
		c.JSONOK("404", nil)
	}))
	// ping server
	g.GET("/ping", ginx.Handle(func(c *ginx.Context) {
		//
	}))
}
