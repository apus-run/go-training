package repository

import (
	"context"
	"database/sql"
	"fmt"
	"project-layout/internal/repository/dao/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"project-layout/internal/domain/entity"
	cachemocks "project-layout/internal/repository/cache/mocks"
	"project-layout/internal/repository/cache/user"
	"project-layout/internal/repository/dao"
	daomocks "project-layout/internal/repository/dao/mocks"
)

func Test_userRepository_FindByID(t *testing.T) {
	// 因为存储的是毫秒数，也就是纳秒部分被去掉了
	// 所以我们需要利用 nowMs 来重建一个不含纳秒部分的 time.Time
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)

	testCases := []struct {
		name string

		// 返回 mock 的 UserDAO 和 UserCache
		mock func(ctrl *gomock.Controller) (dao.UserDAO, user.UserCache)

		// 输入
		ctx context.Context
		id  uint64

		// 预期输出
		wantUser *entity.User
		wantErr  error
	}{
		{
			name: "缓存未命中, 找到了用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, user.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 注意这边，我们传入的是 uint64，
				// 所以要做一个显式的转换，不然默认 12 是 int 类型
				c.EXPECT().GetObj(gomock.Any(), uint64(12)).
					// 缓存未命中
					Return(entity.User{}, user.ErrKeyNotExist)

				// 模拟回写缓存
				c.EXPECT().SetOjb(gomock.Any(), entity.User{
					ID:          12,
					Name:        "小芳",
					Email:       "123@qq.com",
					Password:    "123456",
					Phone:       "13801234567",
					CreatedTime: now,
				}).Return(nil)

				// 查找数据库
				d.EXPECT().FindByID(gomock.Any(), uint64(12)).
					Return(
						&model.User{
							ID:   12,
							Name: "小芳",
							Email: sql.NullString{
								String: "123@qq.com",
								Valid:  true,
							},
							Password: "123456",
							Phone: sql.NullString{
								String: "13801234567",
								Valid:  true,
							},
							CreatedTime: nowMs,
							UpdatedTime: nowMs,
						}, nil)

				return d, c
			},

			ctx: context.Background(),
			id:  12,

			wantUser: &entity.User{
				ID:          12,
				Name:        "小芳",
				Email:       "123@qq.com",
				Password:    "123456",
				Phone:       "13801234567",
				CreatedTime: now,
			},
		},
		{
			name: "缓存命中, 找到了用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, user.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 注意这边，我们传入的是 uint64，
				// 所以要做一个显式的转换，不然默认 12 是 int 类型
				c.EXPECT().GetObj(gomock.Any(), uint64(12)).
					// 模拟缓存命中
					Return(
						entity.User{
							ID:          12,
							Name:        "小芳",
							Email:       "123@qq.com",
							Password:    "123456",
							Phone:       "13801234567",
							CreatedTime: now,
						}, nil)

				return d, c
			},

			ctx: context.Background(),
			id:  12,

			wantUser: &entity.User{
				ID:          12,
				Name:        "小芳",
				Email:       "123@qq.com",
				Password:    "123456",
				Phone:       "13801234567",
				CreatedTime: now,
			},
		},
		{
			name: "没找到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, user.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 注意这边，我们传入的是 uint64，
				// 所以要做一个显式的转换，不然默认 12 是 int 类型
				c.EXPECT().GetObj(gomock.Any(), uint64(12)).
					// 模拟缓存未命中
					Return(entity.User{}, user.ErrKeyNotExist)

				// 未找到用户
				d.EXPECT().FindByID(gomock.Any(), uint64(12)).
					Return(&model.User{}, dao.ErrRecordNotFound)

				return d, c
			},

			ctx: context.Background(),
			id:  12,

			wantErr: fmt.Errorf("此用户不存在, err: %v", ErrUserDataNotFound),

			wantUser: &entity.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDao, userCache := tc.mock(ctrl)
			repo := NewUserRepository(userDao, userCache, nil)
			u, err := repo.FindByID(tc.ctx, tc.id)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}
