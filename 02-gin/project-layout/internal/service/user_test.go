package service

import (
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"project-layout/internal/domain/entity"
	"project-layout/internal/repository"
	repomocks "project-layout/internal/repository/mocks"
)

func Test_userService_Login(t *testing.T) {
	// 固定时间
	now := time.Now()
	testCases := []struct {
		// 用例名称及描述
		name string

		// 预期输入, 根据你的方法参数、接收器来设计
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		ctx      context.Context
		email    string
		password string

		// 预期输出, 根据你的方法返回值、接收器来设计
		wantErr  error
		wantUser entity.User
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), gomock.Any()).
					Return(
						&entity.User{
							ID:          1,
							Name:        "moocss",
							Email:       "moocss@160.com",
							Password:    "$2a$10$3FgSVRUpi.9.UferU0SxC.nGw8Y2eh6yjSqGO8nUcF0CfYbKlrZbS",
							Phone:       "13401234567",
							CreatedTime: now,
						}, nil)

				return repo
			},
			ctx:   context.Background(),
			email: "moocss@163.com",
			// 原始密码
			password: "hello#world123",

			wantErr: nil,
			wantUser: entity.User{
				ID:          1,
				Name:        "moocss",
				Email:       "moocss@160.com",
				Password:    "$2a$10$3FgSVRUpi.9.UferU0SxC.nGw8Y2eh6yjSqGO8nUcF0CfYbKlrZbS",
				Phone:       "13401234567",
				CreatedTime: now,
			},
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), gomock.Any()).
					// 返回错误, 模拟返回 ErrUserDataNotFound 错误
					Return(nil, repository.ErrUserDataNotFound)

				return repo
			},
			ctx:   context.Background(),
			email: "moocss@126.com",
			// 原始密码
			password: "hello#world123",
			// 返回密码错误
			wantErr: ErrInvalidUserOrPassword,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), gomock.Any()).
					// 返回错误, 模拟返回 ErrUserDataNotFound 错误
					Return(
						&entity.User{
							ID:          1,
							Name:        "moocss",
							Email:       "moocss@160.com",
							Password:    "$2a$10$3FgSVRUpi.9.UferU0SxC.nGw8Y2eh6yjSqGO8nUcF0CfYbKlrZbS",
							Phone:       "13401234567",
							CreatedTime: now,
						}, nil)
				return repo
			},
			ctx:   context.Background(),
			email: "moocss@163.com",
			// 错误的密码
			password: "hello#world1234",
			// 返回密码错误
			wantErr: ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo, nil)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			require.Equal(t, tc.wantErr, err)
			require.Equal(t, tc.wantUser, *user)
		})
	}
}

func TestPasswordEncrypt(t *testing.T) {
	pwd := []byte("hello#world123")
	// 加密
	encrypted, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	// 比较
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, pwd)
	require.NoError(t, err)
}
