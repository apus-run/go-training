package router

import "C"
import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"gin-with-jwt/model"
	"gin-with-jwt/router/middleware/auth"
)

func Router() http.Handler {
	log.Printf("load router")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(auth.NewBuilder().
		IgnorePaths("/login").
		IgnorePaths("/ping").Build())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		type UserRequest struct {
			Email    string
			Password string
		}

		var req UserRequest
		err := c.Bind(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
			return
		}

		if req.Email == "moocss@163.com" && req.Password == "123456" {
			claims := auth.CustomClaims{}
			claims.UserID = 1
			claims.UserAgent = c.Request.UserAgent()

			// 过期时间为30分钟, 测试使用1分钟过期
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 1))
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := token.SignedString([]byte(auth.TokenKey))
			if err != nil {
				// 生成 token 的字符串失败，算是系统错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "生成 token 失败",
				})
			}
			c.Header("x-jwt-token", tokenStr)

			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "登录成功",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户名或者密码错误",
		})
	})

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")

		data := &model.User{
			Name:  "foo",
			Email: "foo@gmail.com",
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": gin.H{
				"id":    id,
				"name":  data.Name,
				"email": data.Email,
			},
		})
	})

	return r
}
