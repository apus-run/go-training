package dao

import (
	"gorm.io/gorm"

	"gin-with-database/internal/repo/dao/model"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}
