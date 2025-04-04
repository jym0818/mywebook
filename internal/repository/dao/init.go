package dao

import (
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
