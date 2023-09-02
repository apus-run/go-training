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
	"project-layout/internal/web"
	"project-layout/internal/web/handler"
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
		user.NewUserRedisCache,
		// code.NewCodeRedisCache,
		code.NewCodeMemoryCache,

		// Repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		service.NewUserService,
		InitSmsService,
		service.NewCodeService,

		// Handler 部分
		handler.NewUserHandler,
		web.NewRouter,
		InitWebServer,
	))
}
