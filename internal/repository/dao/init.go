package dao

import (
	dao2 "github.com/jym/webook-interactive/repository/dao"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Article{}, &PublishedArticle{}, &dao2.Interactive{})
	return err
}
