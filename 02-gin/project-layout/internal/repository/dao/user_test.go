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
				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO `user` .*").
					WillReturnResult(mockRes)

				assert.NoError(t, err)
				return db
			},
			ctx: context.Background(),
			user: model.User{
				Email: sql.NullString{
					String: "moocss@163.com",
					Valid:  true,
				},
			},
		},
		{
			name: "插入失败-邮箱冲突",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `user` .*").
					WillReturnError(&mysqlDriver.MySQLError{
						Number: 1062,
					})

				assert.NoError(t, err)
				return db
			},
			ctx:     context.Background(),
			user:    model.User{},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "插入失败",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO `user` .*").
					WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx:     context.Background(),
			user:    model.User{},
			wantErr: errors.New("数据库错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
				// SELECT VERSION;
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// 你 mock DB 不需要 ping
				DisableAutomaticPing: true,

				// GORM默认在事务中执行单个的创建、更新、删除操作，以确保数据库数据的完整性
				// ------------------------
				// 理论上让 GORM 执行
				// INSERT XXX
				// ------------------------
				// 实际上 GORM
				// BEGIN;
				// INSERT
				// COMMIT;
				// 测试跳过事务
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
