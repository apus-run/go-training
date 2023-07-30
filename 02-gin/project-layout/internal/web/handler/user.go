package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-layout/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	ctx.String(http.StatusOK, "登录成功")
}

func (h *UserHandler) Register(ctx *gin.Context) {
	ctx.String(http.StatusOK, "注册成功")
}
