//go:build wireinject

package main

import (
	"github.com/google/wire"
	"project-layout/internal/infra"
	"project-layout/internal/repository"
	"project-layout/internal/repository/cache/code"
	"project-layout/internal/repository/cache/user"
	"project-layout/internal/repository/dao"
	"project-layout/internal/service"
	"project-layout/internal/service/oauth2/wechat"
	"project-layout/internal/web/handler"
	"project-layout/internal/web/handler/jwt"
	"project-layout/pkg/ginx"
	"project-layout/pkg/log"
)

// runApp init web application.
func runApp(*log.Logger) (*ginx.HttpServer, func(), error) {
	panic(wire.Build(
		// 数据库 和 缓存
		infra.NewDB,
		infra.NewRDB,
		// infra.NewLocalDB,
		infra.NewData,

		// DAO 部分
		dao.NewUserDAO,

		// Cache 部分
		user.NewUserRedisCache,
		code.NewRedisCodeCache,
		// code.NewCodeMemoryCache,

		// Repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		service.NewUserService,
		InitSmsService,
		service.NewCodeService,
		wechat.NewService,

		// Handler 部分
		jwt.NewJwtHandler,
		handler.NewUserHandler,
		handler.NewOAuth2WechatHandler,

		// gin 中间件
		InitMiddlewares,

		// Web 服务器
		InitWebServer,
	))
}
