package main

import (
	"log"
	"net/http"
	"time"

	redisx "github.com/apus-run/sea-kit/cache/redis"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gin-with-database/repo"
	"gin-with-database/repo/dao"
	"gin-with-database/router"
	"gin-with-database/router/handler"
	"gin-with-database/svc"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:13306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "123456",
		DB:       1,
	})
	return cmd
}

func main() {

	db := InitDB()
	client := InitRedis()
	rdb := redisx.NewCache(client)
	// mdb := memory.NewCache()

	userDAO := dao.NewUserDAO(db)
	userRepository := repo.NewUserRepository(userDAO, rdb)
	userService := svc.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router.Router(userHandler),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("server listen at %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
