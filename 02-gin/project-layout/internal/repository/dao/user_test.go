package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"project-layout/internal/infra"
	"project-layout/internal/repository/dao/model"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_userDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string

		sqlmock func(t *testing.T) *sql.DB

		// 输入
		ctx  context.Context
		user model.User

		wantErr error
	}{
		{
			name: "插入成功",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO `user` .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx:  context.Background(),
			user: model.User{},
		},
		{
			name: "插入失败-邮箱冲突",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO `user` .*").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx:     context.Background(),
			wantErr: ErrUserDuplicate,
		},
		{
			name: "插入失败",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO `user` .*").
					WillReturnError(errors.New("mock db error"))
				return db
			},
			ctx:     context.Background(),
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})

			// 初始化 DB 不能出错，所以这里要断言必须为 nil
			assert.NoError(t, err)

			data := &infra.Data{
				DB: db,
			}

			dao := NewUserDAO(data)
			c, err := dao.Insert(tc.ctx, tc.user)
			t.Logf("插入成功: %v", c)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
