package router

import (
	"context"
	ginslog "gin-with-middleware/router/middleware/slog"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gin-with-middleware/dao/model"
	"gin-with-middleware/router/middleware/accesslog"
	"gin-with-middleware/router/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Router() http.Handler {
	log.Printf("load web")

	// Create a slog logger, which:
	//   - Logs to stdout.
	w := os.Stdout
	logger := slog.New(
		slog.NewJSONHandler(
			w,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)
	logger.WithGroup("http").
		With("environment", "production").
		With("server", "gin/1.9.0").
		With("server_start_time", time.Now()).
		With("gin_mode", gin.EnvGinMode)
	// [SetDefault]还更新了[log]包使用的默认logger
	slog.SetDefault(logger)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(auth.NewBuilder().
		IgnorePaths("/login").
		IgnorePaths("/signup").
		IgnorePaths("/ping").Build())
	engine.Use(ginslog.NewBuilder(logger).Build())
	engine.Use(accesslog.NewBuilder(
		func(ctx context.Context, al accesslog.AccessLog) {
			logger.Debug("Gin 收到请求", slog.Any("req", al))
		}).AllowReqBody().AllowRespBody().Build())

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	engine.POST("/login", func(c *gin.Context) {
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
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
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

	engine.POST("/signup", func(c *gin.Context) {
		type SignupReq struct {
			Name            string `json:"name"`
			Email           string `json:"email"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirmPassword"`
		}

		var req SignupReq
		// 当我们调用 Bind 方法的时候，如果有问题，Bind 方法已经直接写响应回去了
		if err := c.Bind(&req); err != nil {
			return
		}

		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "用户名不能为空",
			})
			return
		}
		if len(req.Name) < 2 || len(req.Name) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "真实姓名为2 ~ 20个字符",
			})
			return
		}

		if len(req.Email) < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "邮箱不能为空",
			})

			return
		}

		if len(req.Password) < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "密码不能为空",
			})

			return
		}

		if len(req.ConfirmPassword) < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "确认密码不能为空",
			})

			return
		}

		if req.Password != req.ConfirmPassword {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "两次输入的密码不一致",
			})
			return
		}

		data := &model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.ConfirmPassword,
		}

		log.Println(data)

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "注册成功",
		})

	})

	engine.GET("/user/:id", func(c *gin.Context) {
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

	engine.GET("/profile", func(c *gin.Context) {
		data := &model.User{
			Name:  "foo",
			Email: "foo@gmail.com",
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": gin.H{
				"name":  data.Name,
				"email": data.Email,
			},
		})
	})

	return engine
}
