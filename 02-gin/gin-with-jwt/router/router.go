package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-with-jwt/model"
)

func Handler() http.Handler {
	log.Printf("load router")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	r.GET("/me", func(c *gin.Context) {
		u := &model.User{
			Name:  "moocss",
			Email: "moocss@163,com",
		}
		c.JSON(http.StatusOK, gin.H{
			"user": u,
		})

		return
	})

	return r
}
