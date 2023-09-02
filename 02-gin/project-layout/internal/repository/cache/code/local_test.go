package code

import (
	"context"
	"fmt"
	"github.com/coocood/freecache"
	"testing"
	"time"
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
	cache := NewCodeMemoryCache()
	err := cache.Set(context.Background(), "login", "13812345678", "12345")
	if err != nil {
		return
	}

	// 等待 10 秒
	time.Sleep(10 * time.Second)

	ok, err := cache.Verify(context.Background(), "login", "13812345678", "1234")
	if err != nil {
		t.Logf("验证 %v", err)
	}
	t.Logf("验证结果 %v", ok)

}
