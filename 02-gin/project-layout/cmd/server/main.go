package main

import (
	"flag"

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
	c := conf.New(
		conf.WithSource(
			file.NewSource(flagconf),
		),
	)
	c.Load()
	c.Watch()

	logger := log.NewZapLogger(
		log.WithEncoding("json"),
		log.WithFilename("../logs/server.log"),
	)

	defer logger.Sync()

	s, cleanup, err := runApp(logger)
	if err != nil {
		panic(err)
	}
	s.RegisterOnShutdown(cleanup)
}
