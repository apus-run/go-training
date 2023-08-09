package handler

import (
	"net/http"
	"project-layout/internal/domain/entity"
	"project-layout/internal/web/dto"
	"project-layout/pkg/jwt"
	"project-layout/pkg/log"
	"time"

	"project-layout/internal/service"
	"project-layout/pkg/ginx"
)

type UserHandler struct {
	svc service.UserService

	log *log.Logger
}

func NewUserHandler(svc service.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		svc: svc,
		log: logger,
	}
}

func (h *UserHandler) Login(ctx *ginx.Context) {
	var req dto.UserRequest
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSONE(http.StatusOK, "账号或者密码不正确，请重试", nil)
		return
	}

	// 过期时间为30分钟
	expireAt := time.Now().Add(time.Minute * 30)
	token, err := jwt.GenerateToken(
		jwt.WithSecretKey("secret"),
		jwt.WithUserId(user.ID),
		jwt.WithExpireAt(expireAt),
	)

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSONOK("登录成功", struct {
		Token    string    `json:"token"`
		ExpireAt time.Time `json:"expire_at"`
	}{
		Token:    token,
		ExpireAt: expireAt,
	})
}

func (h *UserHandler) Register(ctx *ginx.Context) {
	var req dto.RegisterRequest
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
	}

	_, err = h.svc.Register(
		ctx.Request.Context(),
		entity.User{
			Name:     req.Name,
			Phone:    req.Phone,
			Email:    req.Email,
			Password: req.ConfirmPassword,
		},
	)

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSONOK("注册成功", nil)
}
