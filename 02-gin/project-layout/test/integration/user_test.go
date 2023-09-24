//go:build e2e

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"project-layout/internal/infra"
	"project-layout/pkg/ginx"
	"testing"
	"time"
)

func TestUserHandler_e2e_Login(t *testing.T) {
	const sendSMSCodeUrl = "/v1/login_sms/code/send"
	server := Start()
	rdb := infra.NewRDB()

	testCases := []struct {
		name string

		// 提前准备数据
		before func(t *testing.T)

		// 验证并且删除数据
		after func(t *testing.T)

		phone string

		// 预期响应
		wantCode   int
		wantResult ginx.Response
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 啥也不做
			},
			// 在设置成功的情况下，我们预期在 Redis 里面会有这个数据
			after: func(t *testing.T) {
				ctx := context.Background()
				key := "phone_code:login:15212345678"
				val, err := rdb.Get(ctx, key).Result()
				// 断言必然取到了数据
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
				// 这里可以考虑进一步断言过期时间
				ttl, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				// 过期时间是十分钟，所以这里肯定会大于 9 分钟
				assert.True(t, ttl > time.Minute*9)

				// 删除数据
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "15212345678",
			wantCode: 200,
			wantResult: ginx.Response{
				Code: 200,
				Msg:  "发送成功",
				Data: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			body := fmt.Sprintf(`{"phone": "%s"}`, tc.phone)
			req, err := http.NewRequest(http.MethodPost, sendSMSCodeUrl,
				bytes.NewBuffer([]byte(body)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			code := recorder.Code
			// 反序列化为结果
			var result ginx.Response
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, code)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}
