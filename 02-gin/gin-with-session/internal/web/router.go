package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-with-seesion/internal/web/handler"
	"gin-with-seesion/internal/web/middleware/auth"
)

func Router(uh *handler.UserHandler) http.Handler {
	log.Printf("load web")

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(
		auth.NewBuilder().
			IgnorePaths("/login").
			IgnorePaths("/signup").
			IgnorePaths("/ping").Build(),
		gin.Recovery(),
	)

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	engine.POST("/login", uh.Login)
	engine.POST("/signup", uh.Signup)
	engine.GET("/profile", uh.Profile)

	return engine
}
