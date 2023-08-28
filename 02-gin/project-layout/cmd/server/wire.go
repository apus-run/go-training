//go:build wireinject

package main

import (
	"project-layout/pkg/ginx"
	"project-layout/pkg/log"
)

// wireApp init web application.
func wireApp(*log.Logger) (*ginx.HttpServer, func(), error) {
	panic(wire.Build(
		// 数据库 和 缓存
		infra.NewDB,
		infra.NewRDB,
		infra.NewData,

		// DAO 部分
		dao.NewUserDAO,

		// Cache 部分
		cache.NewUserCache,
		cache.NewCodeCache,

		// Repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		service.NewUserService,
		sms.InitSmsService,
		service.NewCodeService,

		// Handler 部分
		handler.NewUserHandler,
		web.NewRouter,
		InitWebServer,
	))
}
