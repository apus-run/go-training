package infra

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
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
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)
	if err != nil {
		panic(err)
	}

	// 没有开启debug模式，不打印sql
	// db.Logger = logger.Default.LogMode(logger.Silent)

	// 为了方便，我们这里直接把表初始化放在这里
	model.InitTables(db)

	return db
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
