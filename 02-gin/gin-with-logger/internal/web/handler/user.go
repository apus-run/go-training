package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"gin-with-logger/internal/domain"
	"gin-with-logger/internal/web/middleware/auth"
)

type UserHandler struct {
	log *slog.Logger
}

func NewUserHandler(logger *slog.Logger) *UserHandler {
	return &UserHandler{
		log: logger,
	}
}

func (uh *UserHandler) Login(c *gin.Context) {
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

}

func (uh *UserHandler) Signup(c *gin.Context) {
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

	data := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.ConfirmPassword,
	}

	uh.log.Info("注册成功:", data)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

func (uh *UserHandler) Profile(c *gin.Context) {
	type Profile struct {
		Email string
		Phone string
	}

	c.JSON(http.StatusOK, Profile{
		Email: "moocss@126.com",
		Phone: "13401234567",
	})
}
