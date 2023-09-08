package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/apus-run/sea-kit/cache/memory"
	"github.com/apus-run/sea-kit/config"
	"github.com/apus-run/sea-kit/config/file"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gin-with-database/internal/conf"
	"gin-with-database/internal/repo"
	"gin-with-database/internal/repo/dao"
	"gin-with-database/internal/svc"
	"gin-with-database/internal/web"
	"gin-with-database/internal/web/handler"
)

func InitDB(conf conf.Conf) *gorm.DB {
	db, err := gorm.Open(mysql.Open(conf.Data.Database.Dsn))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitRedis(conf conf.Conf) redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Data.Redis.Addr,
		Password: conf.Data.Redis.Password,
		DB:       conf.Data.Redis.Db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return rdb
}

func main() {
	path := "./config/config.yaml"
	c := config.New(
		config.WithSource(
			// 添加前缀为 WEBSERVER_ 的环境变量，不需要的话也可以设为空字符串
			// env.NewSource("WEBSERVER_"),
			file.NewSource(path),
		),
	)

	defer c.Close()

	// 加载配置源：
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("配置文件: %+v", c)
	//
	//dsn, err := c.Value("data.database.dsn").String()
	//if err != nil {
	//	log.Printf("没找到配置字段: %v", err)
	//}
	//log.Printf("data.database.dsn: %v", dsn)

	var cf conf.Conf
	if err := c.Scan(&cf); err != nil {
		panic(err)
	}
	log.Printf("配置文件: %+v", cf)

	db := InitDB(cf)
	// client := InitRedis(cf)
	// rdb := redisx.NewCache(client)
	mdb := memory.NewCache()

	userDAO := dao.NewUserDAO(db)
	userRepository := repo.NewUserRepository(userDAO, mdb)
	userService := svc.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: web.Router(userHandler),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("server listen at %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
