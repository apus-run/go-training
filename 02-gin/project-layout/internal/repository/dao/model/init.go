package model

import (
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}
