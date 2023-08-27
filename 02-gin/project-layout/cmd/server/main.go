package main

import (
	"flag"
	_ "go.uber.org/automaxprocs"
	"project-layout/internal/repository/cache/code"
	"project-layout/internal/service/sms"

	"project-layout/internal/infra"
	"project-layout/internal/repository"
	"project-layout/internal/repository/cache"
	"project-layout/internal/repository/dao"
	"project-layout/internal/service"
	"project-layout/internal/web"
	"project-layout/internal/web/handler"
	"project-layout/pkg/conf"
	"project-layout/pkg/conf/file"
	"project-layout/pkg/ginx"
	"project-layout/pkg/log"
)

// flagconf is the config file path
var flagconf string

func init() {
	// 设置命令行参数
	flag.StringVar(&flagconf, "conf", conf.GetEnvString("CONFIG_PATH", "./config"), "config path, eg: -conf config.yaml")
}

func newApp(log *log.Logger) *ginx.HttpServer {
	// 创建HTTPServer
	srv := ginx.NewHttpServer(
		ginx.WithPort("9000"),
		ginx.WithMode("prod"),
	)

	return srv
}

func main() {
	// 项目配置文件
	flag.Parse()
	c := conf.New(conf.WithSource(
		file.NewSource(flagconf),
	))
	c.Load()
	c.Watch()

	logger := log.NewLogger(
		log.WithEncoding("json"),
		log.WithFilename("../logs/server.log"),
	)

	// 核心业务逻辑
	db := infra.NewDB()
	rdb := infra.NewRDB()
	data, cleanup := infra.NewData(db, rdb, logger)
	ud := dao.NewUserDAO(data)
	uc := cache.NewUserCache(data)
	cc := code.NewCodeCache(rdb)

	ur := repository.NewUserRepository(ud, uc, logger)
	us := service.NewUserService(ur, logger)

	// 方便测试, 本地短信服务模拟短信服务
	// smsMemoryService := localsms.NewService()
	//
	smsSvc := sms.InitSmsService()
	cr := repository.NewCodeRepository(cc)
	cs := service.NewCodeService(smsSvc, cr)

	uh := handler.NewUserHandler(us, cs, logger)

	// 创建HTTPServer
	srv := ginx.NewHttpServer(
		ginx.WithPort("9000"),
		ginx.WithMode("prod"),
	)

	// 关闭 HTTPServer 时清理资源
	srv.RegisterOnShutdown(func() {
		// 释放数据库资源
		cleanup()
	})

	// 注册路由
	r := web.NewRouter(uh)
	srv.Run(r)
}
