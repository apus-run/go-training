package web

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"project-layout/internal/web/handler"
	"project-layout/internal/web/middleware"
	"project-layout/pkg/ginx"
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
	ug := g.Group("/v1/user")
	{
		ug.POST("/login", ginx.Handle(r.userHandler.Login))
		ug.POST("/register", ginx.Handle(r.userHandler.Register))
	}

	ag := g.Group("/v1/article", middleware.AuthToken())
	{
		ag.GET("/list", ginx.Handle(func(c *ginx.Context) {
			fmt.Println("get article list")
		}))
	}
}
