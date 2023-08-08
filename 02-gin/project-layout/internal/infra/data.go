package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"project-layout/internal/repository/dao/model"
	"project-layout/pkg/conf"
	"project-layout/pkg/log"
)

type Data struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// NewData return a new Data
func NewData(db *gorm.DB, rdb *redis.Client, log *log.Logger) (*Data, func()) {
	data := &Data{
		DB:  db,
		RDB: rdb,
	}

	cleanup := func() {
		log.Info("closing the data resources")
		m, err := data.DB.DB()
		if err != nil {
			log.Error(err.Error())
		}
		if err := m.Close(); err != nil {
			log.Error(err.Error())
		}

		if err := data.RDB.Close(); err != nil {
			log.Error(err.Error())
		}
	}

	return data, cleanup
}

func NewDB() *gorm.DB {
	dsn := conf.Get("config", "data.database.dsn").(string)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate
	InitDB(db)

	return db
}

func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.User{},
	); err != nil {
		panic(err)
	}
}

func NewRDB() *redis.Client {
	addr := conf.Get("config", "data.redis.addr").(string)
	pwd := conf.Get("config", "data.redis.password").(string)
	dbname := conf.Get("config", "data.redis.db").(int)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       dbname,
	})

	return rdb
}
