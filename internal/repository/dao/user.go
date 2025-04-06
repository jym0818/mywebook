package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrUserDuplicate = errors.New("唯一索引冲突")
var ErrUserNotExists = gorm.ErrRecordNotFound

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type userDAO struct {
	db *gorm.DB
}

func (u *userDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (u *userDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}
func (u *userDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (u *userDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	err := u.db.WithContext(ctx).Create(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func NewuserDAO(db *gorm.DB) UserDAO {
	return &userDAO{db: db}
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string
	Utime    int64
	Ctime    int64
}
