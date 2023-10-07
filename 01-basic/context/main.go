package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	ctx := NewContext(context.Background(), "moocss")
	userID, age, err := getUserProfile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User info response: use time %v -> %v -> %v\n", time.Since(start), userID, age)
}

func getUserProfile(ctx context.Context) (string, int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	type Result struct {
		userID string
		age    int
		err    error
	}

	resultCh := make(chan Result, 1)
	value := FromContext(ctx)
	fmt.Printf("获取值: %v\n", value)

	// 调用 http client 匿名函数
	go func() {
		id, age, err := getWexinUserInfo()
		resultCh <- Result{
			userID: id,
			age:    age,
			err:    err,
		}
	}()

	select {
	// Done
	// 1. time out
	// 2. cancel主动调用
	case <-ctx.Done():
		return "", 0, ctx.Err()
		// return "", 0, fmt.Errorf("自定义错误 ctx.Err(): %v", time.Now())
	case res := <-resultCh:
		return res.userID, res.age, res.err
	}

}

// http client 调用, 例如: 爬虫服务
func getWexinUserInfo() (string, int, error) {
	// time.Sleep(time.Second * 2) 时间如果大于 context.WithTimeout(ctx, time.Second*1) 时间就会报错超时
	time.Sleep(time.Second * 1)
	return "1001", 28, nil
}
