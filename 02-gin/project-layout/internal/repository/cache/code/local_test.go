package code

import (
	"context"
	"fmt"
	"project-layout/internal/repository/cache/localmocks"
	"strconv"
	"testing"
	"time"

	"github.com/coocood/freecache"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewFreeCache(t *testing.T) {
	cacheSize := 10 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)
	key := []byte("abc")
	val := []byte("def")
	expire := 10 // expire in 10 seconds
	cache.Set(key, val, expire)
	got, err := cache.Get(key)
	if err != nil {
		fmt.Println("查找出错了", err)
	}
	fmt.Printf("查找:%s -> %s\n", key, got)

	// 等待 11 秒
	time.Sleep(11 * time.Second)

	got2, err2 := cache.Get([]byte(key))
	if err != nil {
		fmt.Println("2.查找出错了", err2)
	}
	fmt.Printf("2.查找:%s -> %s\n", key, got2)

	// cache.Del(key)
	ttl, err := cache.TTL([]byte(key))
	if err != nil {
		fmt.Println("过期查找出错了", err == freecache.ErrNotFound)
	}
	fmt.Printf("过期时间: %s -> %v\n", key, ttl)

}

func TestNewCodeMemoryCache(t *testing.T) {
	cacheSize := 10 * 1024 * 1024
	fc := freecache.NewCache(cacheSize)
	cache := NewCodeMemoryCache(fc)
	err := cache.Set(context.Background(), "login", "13812345678", "12345")
	if err != nil {
		return
	}

	// 等待 10 秒
	time.Sleep(10 * time.Second)

	ok, err := cache.Verify(context.Background(), "login", "13812345678", "12345")
	if err != nil {
		t.Logf("验证 %v", err)
	}
	t.Logf("验证结果 %v", ok)

}

func Test_codeMemoryCache_Set(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) FreecacheClient

		// 输入
		ctx   context.Context
		biz   string
		phone string
		code  string

		// 预期输出
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) FreecacheClient {
				client := localmocks.NewMockFreecacheClient(ctrl)
				cacheKey := []byte("phone_code:login:13401234567")
				cacheValue := []byte("123456")
				cntKey := string(cacheKey) + ":cnt"

				client.EXPECT().TTL(cacheKey).Return(uint32(0), nil)
				client.EXPECT().Set(cacheKey, cacheValue, 600).
					Return(nil)
				client.EXPECT().Set(cntKey, strconv.Itoa(3), 600).
					Return(nil)
				
				return client
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "13401234567",
			code:  "123456",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCodeMemoryCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
