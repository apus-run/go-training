package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"project-layout/internal/service"
	svcmocks "project-layout/internal/service/mocks"
	"project-layout/pkg/ginx"
)

func TestUserHandler_Register(t *testing.T) {
	const signupUrl = "/v1/user/register"
	testCases := []struct {
		// 用例名称及描述
		name string

		// 预期输入, 根据你的方法参数、接收器来设计
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 因为 request 的构造过程可能很复杂, 所以我们在这里定义一个 Builder
		reqBody func(t *testing.T) *http.Request

		// 预期输出, 根据你的方法返回值、接收器来设计
		wantCode int
		wantBody ginx.Response
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Register(gomock.Any(), gomock.Any()).
					// 注册成功, UserService.Register 返回了 nil
					Return(nil)

				// 在 signup 这个接口里面，并没有用到的 codeSvc，
				// 所以什么不需要准备模拟调用
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello#123456","confirm_password":"hello#123456","email":"moocss@160.com","phone":"13801234567"}`),
				)

				// 准备请求, 构造 HTTP 请求服务
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}

				return req
			},
			wantCode: 200,
			wantBody: ginx.Response{
				Code: 200,
				Msg:  "注册成功",
				Data: map[string]any{},
			},
		},
		{
			name: "非JSON格式的请求体",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 输入的参数不符合要求, 不会执行到下面的验证和业务部分, 直接返回 nil 即可
				return nil, nil
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello#123456","confirm_password":"hello#123456","email":"`),
				)

				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: ginx.Response{
				Code: 400,
				Msg:  "输入的参数格式不正确",
				Data: map[string]any{},
			},
		},
		{
			name: "密码格式不正确",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 输入的参数不符合要求, 不会执行到下面的验证和业务部分, 直接返回 nil 即可
				return nil, nil
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello","confirm_password":"hello","email":"moocss@160.com","phone":"13801234567"}`),
				)

				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: ginx.Response{
				Code: 200,
				Msg:  "密码必须包含数字、特殊字符，并且长度不能小于 8 位",
				Data: map[string]any{},
			},
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 输入的参数不符合要求, 不会执行到下面的验证和业务部分, 直接返回 nil 即可
				return nil, nil
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello#123456","confirm_password":"hello#12345","email":"moocss@160.com","phone":"13801234567"}`),
				)

				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: ginx.Response{
				Code: 200,
				Msg:  "两次输入的密码不一致",
				Data: map[string]any{},
			},
		},
		{
			name: "手机号不正确",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 输入的参数不符合要求, 不会执行到下面的验证和业务部分, 直接返回 nil 即可
				return nil, nil
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello#123456","confirm_password":"hello#123456","email":"moocss@160.com","phone":"1341234567"}`),
				)

				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: ginx.Response{
				Code: 200,
				Msg:  "手机号不正确",
				Data: map[string]any{},
			},
		},
		{
			name: "邮箱不正确",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 输入的参数不符合要求, 不会执行到下面的验证和业务部分, 直接返回 nil 即可
				return nil, nil
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello#123456","confirm_password":"hello#123456","email":"moocss160.com","phone":"13801234567"}`),
				)

				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: ginx.Response{
				Code: 200,
				Msg:  "邮箱不正确",
				Data: map[string]any{},
			},
		},
		{
			name: "邮箱或者手机号已经存在",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Register(gomock.Any(), gomock.Any()).
					// 注册失败, UserService.Register 返回了 ErrUserDuplicate
					Return(service.ErrUserDuplicate)

				// 在 signup 这个接口里面，并没有用到的 codeSvc，
				// 所以什么不需要准备模拟调用
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello#123456","confirm_password":"hello#123456","email":"moocss@160.com","phone":"13801234567"}`),
				)

				// 准备请求, 构造 HTTP 请求服务
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}

				return req
			},
			wantCode: 200,
			wantBody: ginx.Response{
				Code: 200,
				Msg:  "邮箱或者手机号已经存在",
				Data: map[string]any{},
			},
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Register(gomock.Any(), gomock.Any()).
					// 注册失败, 系统本身的异常
					Return(errors.New("服务器异常"))

				// 在 signup 这个接口里面，并没有用到的 codeSvc，
				// 所以什么不需要准备模拟调用
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBody: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer(
					[]byte(`{"name":"admin","password":"hello#123456","confirm_password":"hello#123456","email":"moocss@160.com","phone":"13801234567"}`),
				)

				// 准备请求, 构造 HTTP 请求服务
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}

				return req
			},
			wantCode: 200,
			wantBody: ginx.Response{
				Code: 400,
				Msg:  "服务器异常",
				Data: map[string]any{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			// 利用 mock 来 构造 UserHandler
			uh := NewUserHandler(userSvc, codeSvc, nil)

			// 注册路由
			server := gin.New()
			ug := server.Group("/v1/user")
			{
				ug.POST("/register", ginx.Handle(uh.Register))
			}
			// 准备请求
			req := tc.reqBody(t)
			// 获得 HTTP 响应, 利用 httptest 来记录响应
			// ResponseRecorder is an implementation of http.ResponseWriter that
			// records its mutations for later inspection in tests.
			resp := httptest.NewRecorder()
			// 执行, gin 实现了 ServeHTTP, 调用 ServeHTTP 就是在执行gin
			server.ServeHTTP(resp, req)

			// 断言
			assert.Equal(t, tc.wantCode, resp.Code)
			var result ginx.Response
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, result)
		})
	}
}
