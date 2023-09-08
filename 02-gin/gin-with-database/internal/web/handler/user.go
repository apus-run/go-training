package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"gin-with-database/internal/domain/entity"
	"gin-with-database/internal/svc"
	"gin-with-database/internal/web/middleware/auth"
)

type UserHandler struct {
	svc svc.UserService
}

func NewUserHandler(svc svc.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
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

	user, err := uh.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err == svc.ErrInvalidUserOrPassword {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户名或者密码错误",
		})
		return
	}

	// 测试使用1分钟
	// expireAt := time.Now().Add(time.Minute)
	// 正常设置为30分钟，要将过期时间设置更长一些
	expireAt := time.Now().Add(time.Minute * 30)
	claims := auth.CustomClaims{}
	claims.UserID = user.ID
	claims.UserAgent = c.Request.UserAgent()

	// 过期时间为30分钟, 测试使用1分钟过期
	claims.ExpiresAt = jwt.NewNumericDate(expireAt)
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

}

func (uh *UserHandler) Signup(c *gin.Context) {
	type SignupReq struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
		Phone           string `json:"phone"`
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

	if req.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "手机号不能为空",
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

	err := uh.svc.Register(
		c.Request.Context(),
		entity.User{
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: req.ConfirmPassword,
		},
	)

	if errors.Is(err, svc.ErrUserDuplicate) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "邮箱或者手机号已经存在",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "服务器异常, 注册失败",
		})
		return
	}

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
	cs := c.MustGet("claims").(*auth.CustomClaims)

	log.Println("claims: %v", cs)

	user, err := uh.svc.Profile(c, cs.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户信息没找到",
		})
		return
	}

	c.JSON(http.StatusOK, Profile{
		Email: user.Email,
		Phone: user.Phone,
	})
}
