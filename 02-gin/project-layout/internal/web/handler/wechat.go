package handler

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/lithammer/shortuuid/v4"

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

func (h *OAuth2WechatHandler) OAuth2URL(ctx *ginx.Context) {
	state := uuid.New()
	url, err := h.wechatSvc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSONE(5, "系统错误，请稍后再试", nil)
		return
	}

	ctx.JSONOK("ok", gin.H{
		url: url,
	})
	return
}

func (h *OAuth2WechatHandler) Callback(ctx *ginx.Context) {
	code := ctx.Query("code")
	info, err := h.wechatSvc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSONE(5, "系统错误", nil)
		return
	}
	// 如果查找用户, 如果用户不存在就创建
	_, err = h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSONE(5, "系统错误", nil)
		return
	}

	ctx.JSONOK("登录成功", nil)
}

func (h *OAuth2WechatHandler) Load(engine *gin.Engine) {
	g := engine.Group("/oauth2/wechat")
	g.GET("/:platform/authurl", ginx.Handle(h.OAuth2URL))
	// 这边用 Any 万无一失
	g.Any("/:platform/callback", ginx.Handle(h.Callback))
}
