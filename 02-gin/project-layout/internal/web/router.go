package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"project-layout/internal/web/handler"
	"project-layout/pkg/ginx"
)

type Router struct {
	userHandler *handler.UserHandler
}

func NewRouter(userHandler *handler.UserHandler) ginx.Router {
	router := &Router{
		userHandler: userHandler,
	}
	return router
}

func (r *Router) Load(engine *gin.Engine) {
	// 公共路由
	// -------------------------------------------------------
	// 404
	engine.NoRoute(ginx.Handle(func(c *ginx.Context) {
		c.JSONE(404, "404", nil)
	}))
	// ping server
	engine.GET("/ping", handler.Ping())

	// 用户组
	// -------------------------------------------------------
	ug := engine.Group("/v1/user")
	{
		ug.POST("/login", ginx.Handle(r.userHandler.Login))
		ug.POST("/login_sms/code/send", ginx.Handle(r.userHandler.SendSMSLoginCode))
		ug.POST("/login_sms", ginx.Handle(r.userHandler.LoginSMS))
		ug.POST("/register", ginx.Handle(r.userHandler.Register))
		ug.GET("/profile", ginx.Handle(r.userHandler.Profile))
		ug.POST("/update/profile", ginx.Handle(r.userHandler.UpdateProfile))
	}

	// 新闻组
	// -------------------------------------------------------
	ag := engine.Group("/v1/article")
	{
		ag.GET("/list", ginx.Handle(func(c *ginx.Context) {
			fmt.Println("get article list")
		}))
	}
}
