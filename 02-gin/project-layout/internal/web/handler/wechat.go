package handler

import (
	"github.com/gin-gonic/gin"
	"project-layout/internal/service"
	"project-layout/internal/service/oauth2/wechat"
	ojwt "project-layout/internal/web/handler/jwt"
	"project-layout/pkg/ginx"
)

type OAuth2WechatHandler struct {
	userSvc   service.UserService
	wechatSvc wechat.Service
	ojwt.Handler
}

func NewOAuth2WechatHandler(wechatSvc wechat.Service,
	userSvc service.UserService,
	jwtHdl ojwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		userSvc:   userSvc,
		wechatSvc: wechatSvc,
		Handler:   jwtHdl,
	}
}

func (h *OAuth2WechatHandler) Load(engine *gin.Engine) {
	g := engine.Group("/oauth2/wechat")
	g.GET("/authurl", ginx.Handle(func(c *ginx.Context) {
		c.JSONOK("ok", nil)
	}))
	// 这边用 Any 万无一失
	g.Any("/callback", ginx.Handle(func(c *ginx.Context) {
		c.JSONOK("ok", nil)
	}))
}
