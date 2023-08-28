package main

import (
	"flag"

	_ "go.uber.org/automaxprocs"

	"project-layout/internal/infra"
	"project-layout/internal/repository"
	"project-layout/internal/repository/cache"
	"project-layout/internal/repository/dao"
	"project-layout/internal/service"
	"project-layout/internal/web"
	"project-layout/internal/web/handler"
	"project-layout/pkg/conf"
	"project-layout/pkg/conf/file"
	"project-layout/pkg/log"
)

// flagconf is the config file path
var flagconf string

func init() {
	// 设置命令行参数
	flag.StringVar(&flagconf, "conf", conf.GetEnvString("CONFIG_PATH", "./config"), "config path, eg: -conf config.yaml")
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
	cc := cache.NewCodeCache(rdb)

	ur := repository.NewUserRepository(ud, uc, logger)
	us := service.NewUserService(ur, logger)
	
	smsSvc := InitSmsService()
	cr := repository.NewCodeRepository(cc)
	cs := service.NewCodeService(smsSvc, cr)

	uh := handler.NewUserHandler(us, cs, logger)

	// 注册路由
	r := web.NewRouter(uh)

	// 启动服务
	s := InitWebServer(r)

	// 优雅关闭
	s.RegisterOnShutdown(func() {
		cleanup()
	})
}
