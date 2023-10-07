package wechat

import (
	"context"
	"net/http"

	"project-layout/internal/domain/entity"
	"project-layout/pkg/log"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)

	// VerifyCode 目前大部分公司的 OAuth2 平台都差不多的设计
	// 返回一个 unionId。这个你可以理解为，在第三方平台上的 unionId
	// 你也可以考虑使用 openId 来替换
	// 一家公司如果有很多应用，不同应用都有自建的用户系统
	// 那么 openId 可能更加合适
	VerifyCode(ctx context.Context, code string) (entity.WeChatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client

	logger log.Logger
}

func NewService(appId, appSecret string, logger log.Logger) Service {
	return &service{}
}

func (s service) AuthURL(ctx context.Context, state string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) VerifyCode(ctx context.Context, code string) (entity.WeChatInfo, error) {
	//TODO implement me
	panic("implement me")
}
