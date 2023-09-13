package handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func New() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Login(c *gin.Context) {

}
