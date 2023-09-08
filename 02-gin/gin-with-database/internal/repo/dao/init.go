package dao

import (
	"gin-with-database/internal/repo/dao/model"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}
