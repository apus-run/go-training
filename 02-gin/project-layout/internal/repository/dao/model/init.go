package dao

import (
	"gorm.io/gorm"
	
	"project-layout/internal/repository/dao/model"
)

func InitTables(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}
