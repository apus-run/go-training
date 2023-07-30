package repository

import (
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"project-layout/pkg/conf"
)

type Data struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// NewData return a new Data
func NewData(db *gorm.DB, rdb *redis.Client) (*Data, func(), error) {
	data := &Data{
		DB:  db,
		RDB: rdb,
	}

	cleanup := func() {
		log.Println("closing the data resources")
	}

	return data, cleanup, nil
}

func NewDB() *gorm.DB {
	dsn := conf.Get("config", "data.database.dsn").(string)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func NewRDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your address
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	return rdb
}
