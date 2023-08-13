package web

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"project-layout/internal/web/handler"
	"project-layout/pkg/ginx"
	"project-layout/pkg/ginx/middleware/auth"
	"project-layout/pkg/ginx/middleware/cors"
)

type Router struct {
	userHandler *handler.UserHandler
}

func NewRouter(userHandler *handler.UserHandler) *Router {
	router := &Router{
		userHandler: userHandler,
	}
	return router
}

func (r *Router) Load(g *gin.Engine) {
	// 注册中间件
	// -------------------------------------------------------
	g.Use(cors.NewCORS(
		// 允许前端发送
		cors.WithAllowHeaders([]string{"Content-Type", "Authorization"}),
		// 允许前端获取
		cors.WithExposeHeaders([]string{"x-jwt-token"}),
		cors.WithMaxAge(12*60*60),
	).Build())
	g.Use(auth.NewBuilder().
		IgnorePaths("user/login").
		IgnorePaths("user/register").Build())

	// 公共路由
	// -------------------------------------------------------
	// 404
	g.NoRoute(ginx.Handle(func(c *ginx.Context) {
		c.JSONE(404, "404", nil)
	}))
	// ping server
	g.GET("/ping", handler.Ping())

	// 用户组
	// -------------------------------------------------------
	ug := g.Group("/v1/user")
	{
		ug.POST("/login", ginx.Handle(r.userHandler.Login))
		ug.POST("/register", ginx.Handle(r.userHandler.Register))
		ug.GET("/profile", ginx.Handle(r.userHandler.Profile))
		ug.POST("/update/profile", ginx.Handle(r.userHandler.UpdateProfile))
	}

	// 新闻组
	// -------------------------------------------------------
	ag := g.Group("/v1/article")
	{
		ag.GET("/list", ginx.Handle(func(c *ginx.Context) {
			fmt.Println("get article list")
		}))
	}
}
