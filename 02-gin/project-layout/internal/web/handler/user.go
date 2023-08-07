package handler

import (
	"net/http"
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
	}
	user, err := h.svc.Login(ctx, "username", "password")
	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
	}

	expireAt := time.Now().Add(time.Hour * 24 * 90)
	token, err := jwt.GenerateToken(
		jwt.WithSecretKey("secret"),
		jwt.WithUserId(user.ID),
		jwt.WithExpireAt(expireAt),
	)

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
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
	var req dto.UserRequest
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
	}
	ctx.JSONOK("注册成功", req)
}
