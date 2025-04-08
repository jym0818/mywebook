package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		user    User
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnResult(sqlmock.NewResult(3, 1))

				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "1139022499@qq.com",
					Valid:  true,
				},
				Password: "123456",
			},
			wantErr: nil,
		},
		{
			name: "插入邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnError(&mysql2.MySQLError{
					Number: 1062,
				})

				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "1139022499@qq.com",
					Valid:  true,
				},
				Password: "123456",
			},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "插入失败",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnError(errors.New("test error"))

				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "1139022499@qq.com",
					Valid:  true,
				},
				Password: "123456",
			},
			wantErr: errors.New("test error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			require.NoError(t, err)
			ud := NewuserDAO(db)
			err = ud.Insert(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
